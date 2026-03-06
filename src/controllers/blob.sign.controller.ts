import type { NextFunction, Request, Response } from "express";
import jwt from "jsonwebtoken";
import { z } from "zod";
import { findBlobById } from "#services/blob.service";

const CONFIG = {
    signedUrlTtlMinSeconds: Number(process.env.SIGNED_URL_TTL_MIN_SECONDS ?? 30),
    signedUrlTtlMaxSeconds: Number(process.env.SIGNED_URL_TTL_MAX_SECONDS ?? 900),
};

/**
 * Issues a signed URL payload for private blob download.
 *
 * @route GET /blob/:id/sign
 * @param req Express request (params, query)
 *   - id: string (required, path param)
 *   - ttl: number (optional, query param, seconds)
 * @param res Express response
 * @param next Express error callback
 * @returns 200 OK: Signed URL payload
 * @returns 404 Not Found: Blob not found
 * @returns 410 Gone: Blob expired
 */
export async function getBlobSignedUrl(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        const schema = z.object({
            id: z.string().min(1),
            ttl: z.coerce
                .number()
                .int()
                .min(CONFIG.signedUrlTtlMinSeconds)
                .max(CONFIG.signedUrlTtlMaxSeconds)
                .default(300),
        });
        const parsed = schema.safeParse({ id: req.params.id, ttl: req.query.ttl });
        if (!parsed.success) {
            res
                .status(400)
                .json({ error: "Invalid parameters", details: parsed.error.flatten() });
            return;
        }
        const blob = await findBlobById(parsed.data.id);
        if (!blob) {
            res.status(404).json({ error: "Blob not found" });
            return;
        }
        if (blob.expiresAt && new Date() > new Date(blob.expiresAt)) {
            res.status(410).json({ error: "Blob expired" });
            return;
        }
        let exp = Math.floor(Date.now() / 1000) + parsed.data.ttl;
        if (blob.expiresAt) {
            const expiresAtSec = Math.floor(
                new Date(blob.expiresAt).getTime() / 1000,
            );
            if (expiresAtSec < exp) {
                exp = expiresAtSec;
            }
        }
        const token = jwt.sign(
            { id: blob.id, exp },
            process.env.TOKEN_SECRET ?? "",
        );
        res.json({
            id: blob.id,
            exp,
            token,
            ttl: exp - Math.floor(Date.now() / 1000),
            url: `/blob/${blob.id}?token=${token}`,
        });
        return;
    } catch (error) {
        next(error);
        return;
    }
}
