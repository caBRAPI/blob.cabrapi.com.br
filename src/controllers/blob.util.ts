import crypto from "node:crypto";
import type { Request } from "express";

function getAdminToken(): string | undefined {
  return (
    process.env.TOKEN_SECRET ||
    (process.env.NODE_ENV === "production"
      ? undefined
      : "local-dev-token-secret")
  );
}

function tokenMatches(expectedToken: string, providedToken: string): boolean {
  const expected = Buffer.from(expectedToken);
  const provided = Buffer.from(providedToken);
  if (expected.length !== provided.length) {
    return false;
  }
  return crypto.timingSafeEqual(expected, provided);
}

export function hasAdminAccess(req: Request): boolean {
  const configuredToken = getAdminToken();
  if (!configuredToken) {
    return false;
  }
  const adminTokenHeader = req.headers["x-admin-token"];
  if (typeof adminTokenHeader === "string") {
    return tokenMatches(configuredToken, adminTokenHeader);
  }
  const authorization = req.headers.authorization;
  if (
    typeof authorization === "string" &&
    authorization.startsWith("Bearer ")
  ) {
    return tokenMatches(configuredToken, authorization.slice(7).trim());
  }
  return false;
}
