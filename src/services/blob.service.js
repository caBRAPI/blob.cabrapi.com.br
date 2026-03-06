import fs from "fs/promises"
import db from "#db/db"
import { sha256 } from "#functions/hash.js"
import { buildStoragePath } from "#functions/storagePath.js"

export async function saveBlob(file) {

    const buffer = file.buffer

    const hash = sha256(buffer)

    const path = buildStoragePath(hash)

    await fs.mkdir(path.substring(0, path.lastIndexOf("/")), { recursive: true })

    await fs.writeFile(path, buffer)

    const blob = await db.blob.create({
        data: {
            bucket: "default",
            key: file.originalname,
            filename: file.originalname,
            mime: file.mimetype,
            size: file.size,
            hash,
            path
        }
    })

    return blob
}