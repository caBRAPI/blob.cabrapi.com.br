import { saveBlob } from "./blob.service.js"

export async function upload(req, res) {

    const blob = await saveBlob(req.file)

    res.json(blob)
}