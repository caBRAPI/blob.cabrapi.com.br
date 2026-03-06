import fs from "node:fs/promises";
import path from "node:path";
import type { NextFunction, Request, Response } from "express";
import jwt from "jsonwebtoken";
import { z } from "zod";
import {
    findBlobById,
    incrementBlobDownloadCount,
    resolveBlobAbsolutePath,
} from "#services/blob.service";
import { hasAdminAccess } from "./blob.util";

/**
 * Downloads a blob file if access conditions are satisfied.
 *
 * @route GET /blob/:id
 * @param req Express request (params, query)
 *   - id: string (required, path param)
 *   - token: string (optional, query param for private blobs)
 * @param res Express response
 * @param next Express error callback
 * @returns 200 OK: File stream
 * @returns 403 Forbidden: Invalid or expired signature
 * @returns 404 Not Found: Blob not found
 * @returns 410 Gone: Blob expired
 */
export async function getBlob(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        const schema = z.object({
            id: z.string().min(1),
            token: z.string().optional(),
        });
        const parsed = schema.safeParse({
            id: req.params.id,
            token: req.query.token,
        });
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
        if (!blob.public) {
            if (hasAdminAccess(req)) {
                if (blob.expiresAt && new Date() > new Date(blob.expiresAt)) {
                    res.status(410).json({ error: "Blob expired" });
                    return;
                }
            } else {
                if (!parsed.data.token) {
                    res.status(403).json({
                        error:
                            "Private blob requires signed URL. Use GET /blob/:id/sign first.",
                    });
                    return;
                }
                try {
                    const payload = jwt.verify(
                        parsed.data.token,
                        process.env.TOKEN_SECRET ?? "",
                    ) as { id: string };
                    if (payload.id !== blob.id) throw new Error();
                    if (blob.expiresAt && new Date() > new Date(blob.expiresAt)) {
                        res.status(410).json({ error: "Blob expired" });
                        return;
                    }
                } catch {
                    res.status(403).json({ error: "Invalid or expired signature" });
                    return;
                }
            }
        }
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
