import { defineAppSetup } from "@slidev/types"

import * as storm from "../data/content"

export default defineAppSetup(({ app }) => {
  app.config.globalProperties.$storm = storm
})
