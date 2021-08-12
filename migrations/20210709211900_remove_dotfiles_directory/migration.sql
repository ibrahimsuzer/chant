/*
  Warnings:

  - You are about to drop the `DotfileDir` table. If the table is not empty, all the data it contains will be lost.
  - A unique constraint covering the columns `[path]` on the table `Dotfile` will be added. If there are existing duplicate values, this will fail.

*/
-- DropTable
PRAGMA foreign_keys=off;
DROP TABLE "DotfileDir";
PRAGMA foreign_keys=on;

-- CreateIndex
CREATE UNIQUE INDEX "Dotfile.path_unique" ON "Dotfile"("path");
