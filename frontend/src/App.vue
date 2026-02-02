<template>
  <div class="min-h-screen">
    <header class="px-6 pt-10 pb-6 lg:px-12">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <p class="text-sm font-mono tracking-[0.2em] text-muted-foreground">STORM CONSOLE</p>
          <h1 class="text-4xl font-display font-semibold tracking-tight text-foreground">
            Realtime Messaging Dashboard
          </h1>
          <p class="mt-2 max-w-2xl text-base text-muted-foreground">
            Publish events to NATS and watch them stream live from the gateway SSE endpoint.
          </p>
        </div>
        <div class="flex flex-wrap gap-3">
          <Badge :variant="connected ? 'success' : 'warning'">
            {{ connected ? 'LIVE STREAM' : 'DISCONNECTED' }}
          </Badge>
          <Badge variant="outline">{{ messages.length }} events</Badge>
          <Badge variant="outline">{{ recentRate }} / 10s</Badge>
        </div>
      </div>
    </header>

    <main class="px-6 pb-16 lg:px-12">
      <div class="grid gap-6 lg:grid-cols-[1fr_1.4fr]">
        <Card>
          <CardHeader>
            <div>
              <p class="text-sm font-mono uppercase tracking-[0.2em] text-muted-foreground">Publisher</p>
              <h2 class="text-2xl font-semibold">Send a message</h2>
            </div>
            <div class="flex gap-2">
              <Button variant="ghost" @click="connectStream">Reconnect</Button>
              <Button variant="secondary" @click="disconnectStream">Disconnect</Button>
            </div>
          </CardHeader>
          <CardContent>
            <div class="space-y-4">
              <div>
                <label class="text-sm font-medium text-foreground">Gateway URL</label>
                <Input v-model="gatewayUrl" placeholder="http://localhost:8080" />
              </div>
              <div>
                <label class="text-sm font-medium text-foreground">Subject</label>
                <Input v-model="subject" placeholder="storm.events" />
              </div>
              <div>
                <label class="text-sm font-medium text-foreground">Payload</label>
                <Textarea v-model="payload" rows="6" placeholder='{"type":"hello","msg":"from frontend"}' />
              </div>
              <div class="flex flex-wrap items-center gap-3">
                <Button @click="sendMessage">Publish</Button>
                <Button variant="outline" @click="sendSample">Send sample</Button>
                <span class="text-sm text-muted-foreground">{{ publishStatus }}</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card class="relative overflow-hidden">
          <CardHeader>
            <div>
              <p class="text-sm font-mono uppercase tracking-[0.2em] text-muted-foreground">Live feed</p>
              <h2 class="text-2xl font-semibold">Streaming events</h2>
            </div>
            <div class="text-sm text-muted-foreground">
              Last event: {{ lastEventAt || 'waiting...' }}
            </div>
          </CardHeader>
          <CardContent>
            <div class="h-[420px] overflow-auto rounded-xl border border-border bg-white/70 p-4 shadow-glow">
              <div v-if="messages.length === 0" class="text-sm text-muted-foreground">
                No events yet. Publish a message to see it here.
              </div>
              <div v-else class="space-y-3">
                <div
                  v-for="event in messages"
                  :key="event.id"
                  class="rounded-lg border border-border bg-white px-3 py-2 shadow-sm"
                >
                  <div class="flex items-center justify-between">
                    <span class="text-xs font-mono text-muted-foreground">{{ event.time }}</span>
                    <Badge variant="outline">{{ event.bytes }} bytes</Badge>
                  </div>
                  <pre class="mt-2 whitespace-pre-wrap text-sm text-foreground">{{ event.text }}</pre>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div class="mt-8 grid gap-6 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <p class="text-sm font-mono uppercase tracking-[0.2em] text-muted-foreground">Connection</p>
            <h3 class="text-xl font-semibold">Stream state</h3>
          </CardHeader>
          <CardContent>
            <p class="text-sm text-muted-foreground">{{ connectionHint }}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <p class="text-sm font-mono uppercase tracking-[0.2em] text-muted-foreground">Throughput</p>
            <h3 class="text-xl font-semibold">Recent activity</h3>
          </CardHeader>
          <CardContent>
            <p class="text-sm text-muted-foreground">
              {{ recentRate }} messages in the last 10 seconds
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <p class="text-sm font-mono uppercase tracking-[0.2em] text-muted-foreground">Tip</p>
            <h3 class="text-xl font-semibold">Storm Day drill</h3>
          </CardHeader>
          <CardContent>
            <p class="text-sm text-muted-foreground">
              Open multiple tabs, publish bursts, and watch the stream stay responsive.
            </p>
          </CardContent>
        </Card>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue"
import Badge from "./components/ui/Badge.vue"
import Button from "./components/ui/Button.vue"
import Card from "./components/ui/Card.vue"
import CardContent from "./components/ui/CardContent.vue"
import CardHeader from "./components/ui/CardHeader.vue"
import Input from "./components/ui/Input.vue"
import Textarea from "./components/ui/Textarea.vue"

const gatewayUrl = ref(import.meta.env.VITE_GATEWAY_URL || "http://localhost:8080")
const subject = ref("storm.events")
const payload = ref("{\n  \"type\": \"hello\",\n  \"message\": \"Storm realtime test\"\n}")
const publishStatus = ref("")
const connected = ref(false)
const lastEventAt = ref("")
const messages = ref([])
const recentWindow = ref([])
let stream = null
let nextId = 1
let intervalId = null

const connectionHint = computed(() =>
  connected.value
    ? "Streaming events over SSE from the gateway."
    : "Stream disconnected. Check gateway URL or reconnect."
)

const recentRate = computed(() => recentWindow.value.length)

const pushEvent = (text) => {
  const now = new Date()
  const entry = {
    id: nextId++,
    time: now.toLocaleTimeString(),
    text,
    bytes: new TextEncoder().encode(text).length,
  }
  messages.value.unshift(entry)
  if (messages.value.length > 120) {
    messages.value.pop()
  }
  lastEventAt.value = now.toLocaleTimeString()
  recentWindow.value.push(now.getTime())
  pruneWindow()
}

const pruneWindow = () => {
  const cutoff = Date.now() - 10000
  recentWindow.value = recentWindow.value.filter((t) => t >= cutoff)
}

const connectStream = () => {
  disconnectStream()
  const url = `${gatewayUrl.value}/events?subject=${encodeURIComponent(subject.value)}`
  stream = new EventSource(url)
  connected.value = true

  stream.onmessage = (event) => {
    if (event?.data) {
      pushEvent(event.data)
    }
  }

  stream.onerror = () => {
    connected.value = false
  }
}

const disconnectStream = () => {
  if (stream) {
    stream.close()
    stream = null
  }
  connected.value = false
}

const sendMessage = async () => {
  publishStatus.value = "publishing..."
  try {
    const res = await fetch(
      `${gatewayUrl.value}/publish?subject=${encodeURIComponent(subject.value)}`,
      {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: payload.value,
      }
    )
    if (!res.ok) {
      throw new Error(await res.text())
    }
    publishStatus.value = "published"
  } catch (err) {
    publishStatus.value = `failed: ${err.message}`
  } finally {
    setTimeout(() => {
      publishStatus.value = ""
    }, 2000)
  }
}

const sendSample = () => {
  payload.value = JSON.stringify(
    {
      type: "sample",
      ts: new Date().toISOString(),
      source: "frontend",
    },
    null,
    2
  )
  sendMessage()
}

onMounted(() => {
  connectStream()
  intervalId = setInterval(pruneWindow, 1000)
})

onBeforeUnmount(() => {
  disconnectStream()
  if (intervalId) {
    clearInterval(intervalId)
  }
})
</script>
