import express from "express"
import multer from "multer"
import { upload } from "#controllers/blob.controller.js"

const router = express.Router()

const uploadMiddleware = multer()

router.post("/upload", uploadMiddleware.single("file"), upload)

export default router