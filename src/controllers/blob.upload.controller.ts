import type { NextFunction, Request, Response } from "express";
import { z } from "zod";
import { saveBlob } from "#services/blob.service";
import { hasAdminAccess } from "./blob.util";

/**
 * Uploads a new blob (file) to storage.
 *
 * @route POST /blob/upload
 * @param req Express request (multipart/form-data)
 *   - file: File (required)
 *   - bucket: string (optional)
 *   - key: string (optional)
 *   - public: boolean (optional)
 *   - metadata: JSON string (optional)
 *   - expiresAt: ISO string or timestamp (optional)
 * @param res Express response
 * @param next Express error callback
 * @returns 201 Created: Blob metadata
 * @returns 400 Bad Request: Missing file or invalid input
 * @returns 401 Unauthorized: Missing/invalid admin token
 */
export async function uploadBlob(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        if (!hasAdminAccess(req)) {
            res
                .status(401)
                .json({ error: "Administrative token required for upload" });
            return;
        }
        if (!req.file) {
            res.status(400).json({ error: "Missing file field in multipart body" });
            return;
        }
        const schema = z.object({
            bucket: z.string().optional(),
            key: z.string().optional(),
            public: z
                .preprocess((v) => v === "true" || v === true, z.boolean())
                .optional(),
            metadata: z.string().optional(),
            expiresAt: z.string().datetime().optional(),
        });
        const parsed = schema.safeParse(req.body);
        if (!parsed.success) {
            res
                .status(400)
                .json({ error: "Invalid input", details: parsed.error.flatten() });
            return;
        }
        let expiresAt: Date | undefined;
        if (parsed.data.expiresAt) {
            const date = new Date(parsed.data.expiresAt);
            if (!Number.isNaN(date.getTime())) {
                expiresAt = date;
            }
        }
        const blob = await saveBlob(req.file, {
            bucket: parsed.data.bucket,
            key: parsed.data.key,
            isPublic: parsed.data.public,
            metadata: parsed.data.metadata,
            expiresAt,
        });
        res.status(201).json(blob);
        return;
    } catch (error) {
        next(error);
        return;
    }
}
