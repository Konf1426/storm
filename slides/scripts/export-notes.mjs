import { mkdirSync, readFileSync, writeFileSync } from "node:fs"
import { dirname, resolve } from "node:path"

const sourcePath = resolve(process.cwd(), "slides.md")
const outputPath = resolve(process.cwd(), "exports", "storm-defense-notes.md")

const source = readFileSync(sourcePath, "utf8").replace(/\r\n/g, "\n")

function stripFrontmatter(input) {
  if (!input.startsWith("---\n")) {
    return input
  }

  const end = input.indexOf("\n---\n", 4)
  return end === -1 ? input : input.slice(end + 5)
}

function splitSlides(input) {
  return stripFrontmatter(input)
    .split(/\n---\n/g)
    .map(chunk => chunk.trim())
    .filter(Boolean)
}

function stripTags(input) {
  return input
    .replace(/<[^>]+>/g, " ")
    .replace(/\s+/g, " ")
    .trim()
}

function extractTitle(slide, index) {
  const heading = slide.match(/^#{1,3}\s+(.+)$/m)
  if (heading) {
    return heading[1].trim()
  }

  const htmlHeading = slide.match(/<h[1-3][^>]*>([\s\S]*?)<\/h[1-3]>/i)
  if (htmlHeading) {
    const value = stripTags(htmlHeading[1])
    if (!value.includes("{{")) {
      return value
    }
  }

  const boldLine = slide.match(/^\*\*(.+)\*\*$/m)
  if (boldLine) {
    return boldLine[1].trim()
  }

  return index === 1 ? "STORM" : `Slide ${index}`
}

function extractNotes(slide) {
  const notes = slide.match(/<!--([\s\S]*?)-->\s*$/)
  return notes ? notes[1].trim() : "Aucune note orateur sur cette slide."
}

const slides = splitSlides(source)
const output = [
  "# STORM - Notes orateur",
  "",
  "Document genere automatiquement depuis `slides/slides.md`.",
  "",
  ...slides.flatMap((slide, index) => [
    `## Slide ${index + 1} - ${extractTitle(slide, index + 1)}`,
    "",
    extractNotes(slide),
    "",
  ]),
].join("\n")

mkdirSync(dirname(outputPath), { recursive: true })
writeFileSync(outputPath, output, "utf8")

console.log(`Notes exportees vers ${outputPath}`)
