import type { NextFunction, Request, Response } from "express";
import { z } from "zod";
import { deleteBlobById } from "#services/blob.service";
import { hasAdminAccess } from "#controllers/blob.util";

/**
 * Soft-deletes a blob and removes its file.
 *
 * @route DELETE /blob/:id
 * @summary Delete a blob by ID. Requires admin authentication.
 *
 * @param req Express request (params)
 *   - id: string (required, path param)
 * @param res Express response
 * @param next Express error callback
 *
 * @returns 204 No Content: Blob deleted
 * @returns 401 Unauthorized: Missing/invalid admin token
 * @returns 404 Not Found: Blob not found
 *
 * @example Request
 *   DELETE /blob/abc123
 *   Headers: x-admin-token: <TOKEN>
 *
 * @security AdminToken
 * @see hasAdminAccess
 */
export async function destroyBlob(
    req: Request,
    res: Response,
    next: NextFunction,
): Promise<void> {
    try {
        if (!hasAdminAccess(req)) {
            res
                .status(401)
                .json({ error: "Administrative token required for delete" });
            return;
        }
        const schema = z.object({ id: z.string().min(1) });
        const parsed = schema.safeParse({ id: req.params.id });
        if (!parsed.success) {
            res
                .status(400)
                .json({ error: "Invalid parameters", details: parsed.error.flatten() });
            return;
        }
        const deleted = await deleteBlobById(parsed.data.id);
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
