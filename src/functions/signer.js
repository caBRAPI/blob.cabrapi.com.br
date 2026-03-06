import crypto from "crypto"
import "dotenv/config"

export function sign(id, exp) {

    const payload = `${id}:${exp}`

    return crypto
        .createHmac("sha256", process.env.SECRET)
        .update(payload)
        .digest("hex")
}