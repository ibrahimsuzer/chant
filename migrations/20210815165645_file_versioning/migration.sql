-- CreateTable
CREATE TABLE "FileVersion" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" DATETIME NOT NULL,
    "content" TEXT NOT NULL,
    "hash" TEXT NOT NULL,
    "fileId" TEXT NOT NULL,
    "predecessorId" TEXT,
    FOREIGN KEY ("fileId") REFERENCES "Dotfile" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("predecessorId") REFERENCES "FileVersion" ("id") ON DELETE SET NULL ON UPDATE CASCADE
);

-- RedefineTables
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Dotfile" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" DATETIME NOT NULL,
    "name" TEXT NOT NULL,
    "path" TEXT NOT NULL,
    "extension" TEXT NOT NULL,
    "mimeType" TEXT NOT NULL,
    "language" TEXT NOT NULL,
    "currentId" TEXT,
    FOREIGN KEY ("currentId") REFERENCES "FileVersion" ("id") ON DELETE SET NULL ON UPDATE CASCADE
);
INSERT INTO "new_Dotfile" ("createdAt", "extension", "id", "language", "mimeType", "name", "path", "updatedAt") SELECT "createdAt", "extension", "id", "language", "mimeType", "name", "path", "updatedAt" FROM "Dotfile";
DROP TABLE "Dotfile";
ALTER TABLE "new_Dotfile" RENAME TO "Dotfile";
CREATE UNIQUE INDEX "Dotfile.path_unique" ON "Dotfile"("path");
CREATE UNIQUE INDEX "Dotfile_currentId_unique" ON "Dotfile"("currentId");
PRAGMA foreign_key_check;
PRAGMA foreign_keys=ON;
