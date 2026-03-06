export function buildStoragePath(hash) {

    const a = hash.slice(0, 2)
    const b = hash.slice(2, 4)

    return `storage/objects/${a}/${b}/${hash}`
}