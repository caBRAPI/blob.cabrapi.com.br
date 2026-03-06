import crypto from "node:crypto";
import fs from "node:fs/promises";
import path from "node:path";
import type { NextFunction, Request, Response } from "express";
import { createHttpError } from "#functions/httpError";
import type { RateLimitResult } from "#functions/rateLimit";
import { consumeRateLimit } from "#functions/rateLimit";
import { createNonce, sign, verifySignature } from "#functions/signer";
import {
    deleteBlobById,
    findBlobById,
    incrementBlobDownloadCount,
    listBlobItems,
    resolveBlobAbsolutePath,
    saveBlob,
} from "#services/blob.service";

// Controller layer: HTTP concerns only (input parsing, rate limits, status codes).

/**
 * Reads a positive integer configuration value from environment.
 *
 * @param name Environment variable name.
 * @param fallback Value used when parsing fails.
 * @returns Positive integer configuration value.
 */
function getEnvInt(name: string, fallback: number): number {
    const parsed = Number(process.env[name] ?? fallback);
    return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
}

const CONFIG = {
    signRateLimitMax: getEnvInt("SIGN_RATE_LIMIT_MAX", 20),
    signRateLimitWindowMs: getEnvInt("SIGN_RATE_LIMIT_WINDOW_MS", 60_000),
    downloadRateLimitMax: getEnvInt("DOWNLOAD_RATE_LIMIT_MAX", 120),
    downloadRateLimitWindowMs: getEnvInt("DOWNLOAD_RATE_LIMIT_WINDOW_MS", 60_000),
    signedUrlTtlMinSeconds: getEnvInt("SIGNED_URL_TTL_MIN_SECONDS", 30),
    signedUrlTtlMaxSeconds: getEnvInt("SIGNED_URL_TTL_MAX_SECONDS", 900),
};

/**
 * Extracts best-effort client IP for rate limiting decisions.
 *
 * @param req Express request.
 * @returns Client IP or fallback identifier.
 */
function getClientIp(req: Request): string {
    return req.ip || String(req.headers["x-forwarded-for"] || "unknown");
}

/**
 * Adds rate limit metadata headers to the response.
 *
 * @param res Express response.
 * @param rateLimit Calculated rate-limit result.
 */
function setRateLimitHeaders(res: Response, rateLimit: RateLimitResult): void {
    res.setHeader("x-ratelimit-remaining", String(rateLimit.remaining));
    res.setHeader(
        "x-ratelimit-reset",
        String(Math.floor(rateLimit.resetAt / 1000)),
    );
}

/**
 * Converts common truthy/falsy payloads to boolean.
 *
 * @param value Raw input value.
 * @param fallback Default value when parsing is ambiguous.
 * @returns Parsed boolean.
 */
function parseBoolean(value: unknown, fallback = false): boolean {
    if (value === undefined) {
        return fallback;
    }

    if (typeof value === "boolean") {
        return value;
    }

    if (typeof value === "string") {
        return value.toLowerCase() === "true";
    }

    return fallback;
}

/**
 * Parses optional boolean query params (`true|false`).
 *
 * @param value Raw query value.
 * @returns `true`, `false` or `undefined` when not provided.
 *
 * @throws {HttpError} When query value is present but invalid.
 */
function parseOptionalBooleanQuery(value: unknown): boolean | undefined {
    if (value === undefined) {
        return undefined;
    }

    if (typeof value === "boolean") {
        return value;
    }

    if (typeof value !== "string") {
        throw createHttpError("Invalid boolean query value", 400);
    }

    const normalized = value.trim().toLowerCase();
    if (normalized === "true") {
        return true;
    }
    if (normalized === "false") {
        return false;
    }

    throw createHttpError("Invalid boolean query value. Use true or false", 400);
}

/**
 * Resolves administrative token from environment.
 *
 * Uses TOKEN_SECRET to keep a single credential for protected operations.
 * Falls back to local dev token in non-production environments.
 */
function getAdminToken(): string | undefined {
    return (
        process.env.TOKEN_SECRET ||
        (process.env.NODE_ENV === "production"
            ? undefined
            : "local-dev-token-secret")
    );
}

/**
 * Compares two tokens using timing-safe comparison.
 */
function tokenMatches(expectedToken: string, providedToken: string): boolean {
    const expected = Buffer.from(expectedToken);
    const provided = Buffer.from(providedToken);

    if (expected.length !== provided.length) {
        return false;
    }

    return crypto.timingSafeEqual(expected, provided);
}

/**
 * Checks whether caller has administrative access for protected operations.
 *
 * @param req Express request.
 * @returns `true` when request contains valid administrative token.
 *
 * Accepted headers:
 * - `x-admin-token: <TOKEN_SECRET>`
 * - `authorization: Bearer <TOKEN_SECRET>`
 */
function hasAdminAccess(req: Request): boolean {
    const configuredToken = getAdminToken();
    if (!configuredToken) {
        return false;
    }

    const adminTokenHeader = req.headers["x-admin-token"];
    if (typeof adminTokenHeader === "string") {
        return tokenMatches(configuredToken, adminTokenHeader);
    }

    const authorization = req.headers.authorization;
    if (
        typeof authorization === "string" &&
        authorization.startsWith("Bearer ")
    ) {
        return tokenMatches(configuredToken, authorization.slice(7).trim());
    }

    return false;
}

/**
 * Handles multipart upload and persists blob metadata/content.
 *
 * Expected payload (`multipart/form-data`):
 * - `file` (required)
 * - `bucket` (optional)
 * - `key` (optional)
 * - `public` (optional, `true|false`)
 * - `metadata` (optional, JSON string)
 *
 * Access control:
 * - requires administrative token (`TOKEN_SECRET`) in header.
 *
 * @param req Express request with multipart body.
 * @param res Express response.
 * @param next Express error callback.
 */
export async function uploadBlob(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        if (!hasAdminAccess(req)) {
            res.status(401).json({
                error: "Administrative token required for upload",
            });
            return;
        }

        if (!req.file) {
            res.status(400).json({ error: "Missing file field in multipart body" });
            return;
        }

        // Permite definir expiresAt opcionalmente (ISO string ou timestamp)
        let expiresAt: Date | undefined = undefined;
        if (req.body.expiresAt) {
            const date = new Date(req.body.expiresAt);
            if (!isNaN(date.getTime())) {
                expiresAt = date;
            }
        }

        const blob = await saveBlob(req.file, {
            bucket: req.body.bucket,
            key: req.body.key,
            isPublic: parseBoolean(req.body.public),
            metadata: req.body.metadata,
            expiresAt,
        });

        res.status(201).json(blob);
        return;
    } catch (error) {
        next(error);
        return;
    }
}

/**
 * Returns paginated blob metadata list.
 *
 * Query params:
 * - `page` (default: 1)
 * - `pageSize` (default: 20, max: 100)
 * - `bucket` (optional)
 * - `public` (optional, `true|false`)
 *
 * Security behavior:
 * - without administrative token, listing is forced to `public=true`
 * - requesting `public=false` without token returns `403`
 *
 * @param req Express request.
 * @param res Express response.
 * @param next Express error callback.
 */
export async function listBlobs(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        const page = Number(req.query.page ?? 1);
        const pageSize = Number(req.query.pageSize ?? 20);
        const bucket =
            typeof req.query.bucket === "string" ? req.query.bucket : undefined;
        const requestedPublic = parseOptionalBooleanQuery(req.query.public);

        // Security default: unauthenticated callers can only list public blobs.
        let isPublic: boolean | undefined = requestedPublic;

        if (!hasAdminAccess(req)) {
            if (requestedPublic === false) {
                res.status(403).json({
                    error: "Listing private blobs requires administrative token",
                });
                return;
            }

            isPublic = true;
        }

        const { data, total } = await listBlobItems({
            page,
            pageSize,
            bucket,
            isPublic,
        });

        res.json({
            page,
            pageSize,
            total,
            items: data,
        });
        return;
    } catch (error) {
        next(error);
        return;
    }
}

/**
 * Streams a blob file if access conditions are satisfied.
 *
 * Private blobs require querystring signature fields: `exp`, `n`, `sig`.
 *
 * @param req Express request.
 * @param res Express response.
 * @param next Express error callback.
 */
export async function getBlob(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        // Apply early throttling to reduce expensive work under abuse.
        const downloadLimit = consumeRateLimit({
            scope: "download",
            key: getClientIp(req),
            max: CONFIG.downloadRateLimitMax,
            windowMs: CONFIG.downloadRateLimitWindowMs,
        });

        setRateLimitHeaders(res, downloadLimit);

        if (!downloadLimit.ok) {
            res.status(429).json({
                error: "Too many requests",
                retryAfterMs: downloadLimit.retryAfterMs,
            });
            return;
        }

        const blobId = String(req.params.id);
        const blob = await findBlobById(blobId);

        if (!blob) {
            res.status(404).json({ error: "Blob not found" });
            return;
        }

        if (!blob.public) {
            // Permite acesso permanente para admin
            if (hasAdminAccess(req)) {
                // Se houver expiresAt e já expirou, bloqueia até para admin
                if (blob.expiresAt && new Date() > new Date(blob.expiresAt)) {
                    res.status(410).json({ error: "Blob expired" });
                    return;
                }
                // Admin pode baixar sem assinatura
            } else {
                // Para acesso externo, exige assinatura válida
                const exp = Number(req.query.exp);
                const sig = req.query.sig;
                const nonce = req.query.n;

                if (!req.query.exp || !req.query.sig || !req.query.n) {
                    res.status(403).json({
                        error:
                            "Private blob requires signed URL. Use GET /blob/:id/sign first.",
                    });
                    return;
                }

                // Se houver expiresAt e já expirou, bloqueia
                if (blob.expiresAt && new Date() > new Date(blob.expiresAt)) {
                    res.status(410).json({ error: "Blob expired" });
                    return;
                }

                if (
                    !verifySignature({
                        id: blob.id,
                        exp,
                        nonce: String(nonce),
                        sig: String(sig),
                        method: req.method,
                    })
                ) {
                    res.status(403).json({ error: "Invalid or expired signature" });
                    return;
                }
            }
        }

        // Async metric update; never block file delivery.
        incrementBlobDownloadCount(blob.id).catch(() => { });

        const absolutePath = resolveBlobAbsolutePath(blob.path);
        await fs.access(absolutePath);

        res.setHeader("Content-Type", blob.mime || "application/octet-stream");
        res.setHeader("Content-Disposition", `inline; filename="${blob.filename}"`);
        res.sendFile(path.resolve(absolutePath));
        return;
    } catch (error) {
        next(error);
        return;
    }
}

/**
 * Soft-deletes blob metadata and removes file when present.
 *
 * Access control:
 * - requires administrative token (`TOKEN_SECRET`) in header.
 *
 * @param req Express request.
 * @param res Express response.
 * @param next Express error callback.
 */
export async function destroyBlob(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        if (!hasAdminAccess(req)) {
            res.status(401).json({
                error: "Administrative token required for delete",
            });
            return;
        }

        const blobId = String(req.params.id);
        const deleted = await deleteBlobById(blobId);

        if (!deleted) {
            res.status(404).json({ error: "Blob not found" });
            return;
        }

        res.status(204).send();
        return;
    } catch (error) {
        next(error);
        return;
    }
}

/**
 * Issues a signed URL payload for private blob download.
 *
 * Query params:
 * - `ttl` in seconds, clamped by environment min/max bounds.
 *
 * @param req Express request.
 * @param res Express response.
 * @param next Express error callback.
 */
export async function getBlobSignedUrl(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        // Limit signed URL minting to reduce brute-force and scraping pressure.
        const signLimit = consumeRateLimit({
            scope: "sign",
            key: `${getClientIp(req)}:${String(req.params.id)}`,
            max: CONFIG.signRateLimitMax,
            windowMs: CONFIG.signRateLimitWindowMs,
        });

        setRateLimitHeaders(res, signLimit);

        if (!signLimit.ok) {
            res.status(429).json({
                error: "Too many signature requests",
                retryAfterMs: signLimit.retryAfterMs,
            });
            return;
        }

        const blobId = String(req.params.id);
        const blob = await findBlobById(blobId);

        if (!blob) {
            res.status(404).json({ error: "Blob not found" });
            return;
        }

        // Se houver expiresAt e já expirou, não gera URL
        if (blob.expiresAt && new Date() > new Date(blob.expiresAt)) {
            res.status(410).json({ error: "Blob expired" });
            return;
        }

        const requestedTtl = Number(req.query.ttl ?? 300);
        const minTtl = CONFIG.signedUrlTtlMinSeconds;
        const maxTtl = CONFIG.signedUrlTtlMaxSeconds;
        const ttlInSeconds = Math.min(maxTtl, Math.max(minTtl, requestedTtl));
        // Se expiresAt for menor que o exp calculado, limita expiração ao expiresAt
        let exp = Math.floor(Date.now() / 1000) + ttlInSeconds;
        if (blob.expiresAt) {
            const expiresAtSec = Math.floor(new Date(blob.expiresAt).getTime() / 1000);
            if (expiresAtSec < exp) {
                exp = expiresAtSec;
            }
        }
        const nonce = createNonce();
        const sig = sign({
            id: blob.id,
            exp,
            nonce,
            method: "GET",
        });

        res.json({
            id: blob.id,
            exp,
            n: nonce,
            sig,
            ttl: exp - Math.floor(Date.now() / 1000),
            url: `/blob/${blob.id}?exp=${exp}&n=${nonce}&sig=${sig}`,
        });
        return;
    } catch (error) {
        next(error);
        return;
    }
}
