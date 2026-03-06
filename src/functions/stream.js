import fs from "fs"

export function streamFile(path, res) {

    const stream = fs.createReadStream(path)

    stream.pipe(res)
}