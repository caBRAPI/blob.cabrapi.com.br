import express from "express";
import multer from "multer";
import { destroyBlob } from "#controllers/blob.delete.controller";
import { getBlob } from "#controllers/blob.get.controller";
import { listBlobs } from "#controllers/blob.list.controller";
import { getBlobSignedUrl } from "#controllers/blob.sign.controller";
import { uploadBlob } from "#controllers/blob.upload.controller";

const router = express.Router();

const maxUploadSize = Number(
    process.env.MAX_UPLOAD_SIZE_BYTES ?? 20 * 1024 * 1024,
);

const uploadMiddleware = multer({
    storage: multer.memoryStorage(),
    limits: {
        fileSize: maxUploadSize,
        files: 1,
    },
});

router.post("/upload", uploadMiddleware.single("file"), uploadBlob);
router.get("/", listBlobs);
router.get("/:id/sign", getBlobSignedUrl);
router.get("/:id", getBlob);
router.delete("/:id", destroyBlob);

export default router;
