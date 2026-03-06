import type { NextFunction, Request, Response } from "express";
import { z } from "zod";
import { listBlobItems } from "#services/blob.service";
import { hasAdminAccess } from "./blob.util";

/**
 * Lists blobs with pagination and optional filters.
 *
 * @route GET /blob
 * @param req Express request (query params)
 *   - page: number (default: 1)
 *   - pageSize: number (default: 20, max: 100)
 *   - bucket: string (optional)
 *   - public: boolean (optional)
 * @param res Express response
 * @param next Express error callback
 * @returns 200 OK: Paginated list of blobs
 * @returns 403 Forbidden: Listing private blobs without admin token
 */
export async function listBlobs(
  req: Request,
  res: Response,
  next: NextFunction,
): Promise<void> {
  try {
    const schema = z.object({
      page: z.coerce.number().int().min(1).default(1),
      pageSize: z.coerce.number().int().min(1).max(100).default(20),
      bucket: z.string().optional(),
      public: z
        .preprocess((v) => v === "true" || v === true, z.boolean())
        .optional(),
    });
    const parsed = schema.safeParse(req.query);
    if (!parsed.success) {
      res
        .status(400)
        .json({ error: "Invalid query", details: parsed.error.flatten() });
      return;
    }
    let isPublic: boolean | undefined = parsed.data.public;
    if (!hasAdminAccess(req)) {
      if (parsed.data.public === false) {
        res.status(403).json({
          error: "Listing private blobs requires administrative token",
        });
        return;
      }
      isPublic = true;
    }
    const { data, total } = await listBlobItems({
      page: parsed.data.page,
      pageSize: parsed.data.pageSize,
      bucket: parsed.data.bucket,
      isPublic,
    });
    res.json({
      page: parsed.data.page,
      pageSize: parsed.data.pageSize,
      total,
      items: data,
    });
    return;
  } catch (error) {
    next(error);
    return;
  }
}
