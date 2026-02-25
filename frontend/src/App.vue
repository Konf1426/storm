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
            Publish events to NATS and watch them stream live over WebSocket.
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
      <div v-if="globalError" class="mt-4 rounded-xl border border-red-200 bg-red-50 p-4 shadow-sm animate-in fade-in slide-in-from-top-2">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 text-sm font-medium text-red-800">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-alert-circle"><circle cx="12" cy="12" r="10"/><line x1="12" x2="12" y1="8" y2="12"/><line x1="12" x2="12.01" y1="16" y2="16"/></svg>
            {{ globalError }}
          </div>
          <Button variant="ghost" size="sm" class="h-8 px-2 text-red-800 hover:bg-red-100" @click="globalError = ''">Dismiss</Button>
        </div>
      </div>
    </header>

    <main class="px-6 pb-16 lg:px-12">
      <div v-if="!authenticated" class="mx-auto max-w-xl">
        <Card>
          <CardHeader>
            <div>
              <p class="text-sm font-mono uppercase tracking-[0.2em] text-muted-foreground">Auth</p>
              <h2 class="text-2xl font-semibold">Sign in to the dashboard</h2>
            </div>
          </CardHeader>
          <CardContent>
            <div class="space-y-4">
              <div>
                <label class="text-sm font-medium text-foreground">User ID</label>
                <Input v-model="loginUser" placeholder="user-1" />
              </div>
              <div>
                <label class="text-sm font-medium text-foreground">Password</label>
                <Input v-model="loginPassword" type="password" placeholder="••••••" />
              </div>
              <div class="flex flex-wrap gap-3">
                <Button @click="login">Login</Button>
                <Button variant="outline" @click="register">Register</Button>
                <span class="text-sm text-muted-foreground">{{ authStatus }}</span>
              </div>
              <p class="text-xs text-muted-foreground">
                Access token lasts 15 min and auto-refreshes with a 24h refresh token.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      <div v-else class="grid gap-6 lg:grid-cols-[1fr_1.4fr]">
        <Card>
          <CardHeader>
            <div>
              <p class="text-sm font-mono uppercase tracking-[0.2em] text-muted-foreground">Publisher</p>
              <h2 class="text-2xl font-semibold">Send a message</h2>
              <p class="mt-1 text-xs text-muted-foreground">Signed in as {{ currentUser }}</p>
            </div>
            <div class="flex gap-2">
              <Button variant="ghost" @click="connectStream">Reconnect</Button>
              <Button variant="secondary" @click="disconnectStream">Disconnect</Button>
              <Button variant="outline" @click="logout">Logout</Button>
            </div>
          </CardHeader>
          <CardContent>
            <div class="space-y-4">
              <div>
                <label class="text-sm font-medium text-foreground">Gateway URL</label>
                <Input v-model="gatewayUrl" placeholder="http://localhost:8080" />
              </div>
              <div v-if="!selectedChannelId">
                <label class="text-sm font-medium text-foreground">Subject</label>
                <Input v-model="subject" placeholder="storm.events" />
                <p class="mt-1 text-xs text-muted-foreground">
                  Used only when no channel is selected.
                </p>
              </div>
              <div>
                <label class="text-sm font-medium text-foreground">Channel</label>
                <div class="mt-2 flex flex-wrap gap-2">
                  <select
                    v-model="selectedChannelId"
                    class="w-full rounded-2xl border border-border bg-white/90 px-4 py-2 text-sm text-foreground shadow-sm focus:border-primary focus:outline-none focus:ring-2 focus:ring-ring/30"
                  >
                    <option value="">No channel (raw subject)</option>
                    <option v-for="channel in channels" :key="channel.id" :value="String(channel.id)">
                      {{ channel.name }} (#{{ channel.id }})
                    </option>
                  </select>
                  <Input v-model="newChannelName" placeholder="new channel name" />
                  <Button variant="outline" @click="createChannel">Create channel</Button>
                </div>
                <div class="mt-2 flex items-center gap-2 text-xs text-muted-foreground">
                  <span>Channels loaded: {{ channels.length }}</span>
                  <Button variant="ghost" @click="loadChannels">Refresh</Button>
                </div>
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
            <div
              ref="feedRef"
              class="h-[420px] overflow-auto rounded-xl border border-border bg-white/70 p-4 shadow-glow"
            >
              <div v-if="messages.length === 0" class="text-sm text-muted-foreground">
                No events yet. Publish a message to see it here.
              </div>
              <div v-else class="space-y-3">
                <div
                  v-for="event in messages"
                  :key="event.id"
                  :class="[
                    'rounded-lg border px-3 py-2 shadow-sm',
                    event.own ? 'border-emerald-200 bg-emerald-50' : 'border-border bg-white'
                  ]"
                >
                  <div class="flex items-center justify-between">
                    <span class="text-xs font-mono text-muted-foreground">
                      {{ event.time }} • {{ event.date }}
                    </span>
                    <Badge variant="outline">{{ event.bytes }} bytes</Badge>
                  </div>
                  <pre class="mt-2 whitespace-pre-wrap text-sm text-foreground">{{ event.text }}</pre>
                </div>
              </div>
            </div>
            <div class="mt-4 flex flex-col gap-3">
              <div>
                <label class="text-sm font-medium text-foreground">Message</label>
                <Input v-model="messageText" placeholder="type your message" />
              </div>
              <div class="flex flex-wrap items-center gap-3">
                <Button @click="sendMessage">Send</Button>
                <span class="text-sm text-muted-foreground">{{ publishStatus }}</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div v-if="authenticated" class="mt-8 grid gap-6 lg:grid-cols-3">
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
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue"
import Badge from "./components/ui/Badge.vue"
import Button from "./components/ui/Button.vue"
import Card from "./components/ui/Card.vue"
import CardContent from "./components/ui/CardContent.vue"
import CardHeader from "./components/ui/CardHeader.vue"
import Input from "./components/ui/Input.vue"

const gatewayUrl = ref(import.meta.env.VITE_GATEWAY_URL || "http://localhost:8080")
const subject = ref("storm.events")
const messageText = ref("")
const publishStatus = ref("")
const authStatus = ref("")
const globalError = ref("")
const connected = ref(false)
const lastEventAt = ref("")
const messages = ref([])
const recentWindow = ref([])
let stream = null
let nextId = 1
let intervalId = null
let refreshTimer = null
let reconnectTimer = null
const feedRef = ref(null)

const channels = ref([])
const selectedChannelId = ref("")
const newChannelName = ref("")

const authenticated = ref(false)
const currentUser = ref("")
const loginUser = ref("")
const loginPassword = ref("")

const connectionHint = computed(() =>
  connected.value
    ? "Streaming events over WebSocket from the gateway."
    : "Stream disconnected. Check gateway URL or reconnect."
)

const recentRate = computed(() => recentWindow.value.length)

const pushEvent = (text) => {
  const parsed = parseMessage(text)
  const now = new Date()
  const entry = {
    id: nextId++,
    time: now.toLocaleTimeString(),
    date: now.toLocaleDateString(),
    text: parsed.text,
    own: parsed.author && parsed.author === currentUser.value,
    bytes: new TextEncoder().encode(text).length,
  }
  const shouldStick = isNearBottom()
  messages.value.push(entry)
  if (messages.value.length > 200) {
    messages.value.shift()
  }
  lastEventAt.value = now.toLocaleTimeString()
  recentWindow.value.push(now.getTime())
  pruneWindow()
  if (shouldStick) {
    scrollToBottom()
  }
}

const pruneWindow = () => {
  const cutoff = Date.now() - 10000
  recentWindow.value = recentWindow.value.filter((t) => t >= cutoff)
}

const connectStream = () => {
  if (!authenticated.value) return
  disconnectStream()
  const base = gatewayUrl.value.replace(/^http/, "ws").replace(/\/$/, "")
  const params = new URLSearchParams()
  if (selectedChannelId.value) {
    params.set("channel_id", selectedChannelId.value)
  } else {
    params.set("subject", subject.value)
  }
  const url = `${base}/ws?${params.toString()}`
  stream = new WebSocket(url)
  stream.onopen = () => {
    connected.value = true
  }
  stream.onmessage = (event) => {
    if (event?.data) {
      pushEvent(event.data)
    }
  }
  stream.onerror = () => {
    connected.value = false
  }
  stream.onclose = () => {
    connected.value = false
    scheduleReconnect()
  }
}

const scheduleReconnect = () => {
  if (!authenticated.value) return
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
  }
  reconnectTimer = setTimeout(() => {
    connectStream()
  }, 1500)
}

const disconnectStream = () => {
  if (stream) {
    stream.close()
    stream = null
  }
  connected.value = false
}

const handleApiError = async (res) => {
  const status = res.status;
  let message = "";
  try {
    message = await res.text();
  } catch (e) {
    message = "Unknown error";
  }

  if (status === 429) {
    globalError.value = "Slow down! You are being rate limited (429). Please wait a moment.";
  } else if (status === 503) {
    globalError.value = "Service temporarily unavailable (503). The backend might be overloaded.";
  } else if (status === 401) {
    globalError.value = "Session expired or invalid (401). Please login again.";
    await logout();
  } else if (status >= 500) {
    globalError.value = `Server error (${status}): ${message || 'Check backend logs'}`;
  } else {
    globalError.value = `Error (${status}): ${message}`;
  }
  
  // Clear error after 10 seconds automatically
  setTimeout(() => {
    if (globalError.value.includes(status.toString())) {
      globalError.value = "";
    }
  }, 10000);
}

const sendMessage = async () => {
  publishStatus.value = "sending..."
  globalError.value = ""
  try {
    let url = `${gatewayUrl.value}/publish?subject=${encodeURIComponent(subject.value)}`
    const payload = {
      user: currentUser.value || "anonymous",
      message: messageText.value.trim(),
    }
    if (!payload.message) {
      publishStatus.value = "message required"
      return
    }
    let body = JSON.stringify(payload)
    if (selectedChannelId.value) {
      url = `${gatewayUrl.value}/channels/${selectedChannelId.value}/messages`
      body = JSON.stringify({ payload: body })
    }
    const res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body,
      credentials: "include",
    })
    
    if (!res.ok) {
      await handleApiError(res);
      throw new Error("Failed to send");
    }
    messageText.value = ""
    publishStatus.value = "sent"
  } catch (err) {
    publishStatus.value = `failed`
  } finally {
    setTimeout(() => {
      publishStatus.value = ""
    }, 2000)
  }
}

const loadChannels = async () => {
  try {
    const res = await fetch(`${gatewayUrl.value}/channels`, {
      credentials: "include",
    })
    if (!res.ok) {
      await handleApiError(res);
      return
    }
    channels.value = await res.json()
  } catch (err) {
    console.error("Failed to load channels", err)
  }
}

const loadHistory = async () => {
  if (!selectedChannelId.value) {
    messages.value = []
    return
  }
  try {
    const res = await fetch(
      `${gatewayUrl.value}/channels/${selectedChannelId.value}/messages?limit=50`,
      { credentials: "include" }
    )
    if (!res.ok) {
      await handleApiError(res);
      return
    }
    const history = await res.json()
    messages.value = history
      .slice()
      .reverse()
      .map((item) => {
        const parsed = parseMessage(item.payload || "", item.user_id)
        return {
          id: `history-${item.id}`,
          time: new Date(item.created_at).toLocaleTimeString(),
          date: new Date(item.created_at).toLocaleDateString(),
          text: parsed.text,
          own: parsed.author && parsed.author === currentUser.value,
          bytes: (item.payload || "").length,
        }
      })
    scrollToBottom()
  } catch (err) {
    console.error("Failed to load history", err)
  }
}

const createChannel = async () => {
  if (!newChannelName.value.trim()) {
    publishStatus.value = "channel name required"
    return
  }
  try {
    const res = await fetch(`${gatewayUrl.value}/channels`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ name: newChannelName.value.trim() }),
    })
    if (!res.ok) {
      await handleApiError(res);
      return
    }
    const channel = await res.json()
    newChannelName.value = ""
    await loadChannels()
    selectedChannelId.value = String(channel.id)
    await loadHistory()
    connectStream()
  } catch (err) {
    console.error("Failed to create channel", err)
  }
}

const register = async () => {
  authStatus.value = "registering..."
  globalError.value = ""
  try {
    const res = await fetch(`${gatewayUrl.value}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        user_id: loginUser.value,
        password: loginPassword.value,
        display_name: loginUser.value,
      }),
    })
    if (!res.ok) {
      await handleApiError(res);
      throw new Error("Register failed")
    }
    authStatus.value = "registered, please login"
  } catch (err) {
    authStatus.value = `register failed`
  }
}

const login = async () => {
  authStatus.value = "logging in..."
  globalError.value = ""
  try {
    const res = await fetch(`${gatewayUrl.value}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        user_id: loginUser.value,
        password: loginPassword.value,
      }),
    })
    if (!res.ok) {
      await handleApiError(res);
      throw new Error("Login failed")
    }
    const user = await res.json()
    authenticated.value = true
    currentUser.value = user.id
    authStatus.value = ""
    await loadChannels()
    await loadHistory()
    connectStream()
    scheduleRefresh()
  } catch (err) {
    authStatus.value = `login failed`
  }
}

const logout = async () => {
  await fetch(`${gatewayUrl.value}/auth/logout`, {
    method: "POST",
    credentials: "include",
  })
  authenticated.value = false
  currentUser.value = ""
  channels.value = []
  messages.value = []
  disconnectStream()
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
}

const refreshSession = async () => {
  try {
    const res = await fetch(`${gatewayUrl.value}/auth/refresh`, {
      method: "POST",
      credentials: "include",
    })
    return res.ok
  } catch (err) {
    console.error("Refresh failed", err)
    return false
  }
}

const scheduleRefresh = () => {
  if (refreshTimer) {
    clearTimeout(refreshTimer)
  }
  // Refresh every 12 minutes (Access Token persists 15 min)
  refreshTimer = setTimeout(async () => {
    if (!authenticated.value) return
    const ok = await refreshSession()
    if (ok) {
      scheduleRefresh()
    } else {
      await logout()
      authStatus.value = "Session expired, please login again"
    }
  }, 12 * 60 * 1000)
}

const checkSession = async () => {
  const res = await fetch(`${gatewayUrl.value}/auth/me`, {
    credentials: "include",
  })
  if (res.ok) {
    const user = await res.json()
    authenticated.value = true
    currentUser.value = user.id
    await loadChannels()
    await loadHistory()
    connectStream()
    scheduleRefresh()
  }
}

const parseMessage = (raw, fallbackAuthor = "") => {
  try {
    const parsed = JSON.parse(raw)
    if (parsed && parsed.user && parsed.message) {
      return {
        text: `${parsed.user} : ${parsed.message}`,
        author: parsed.user,
      }
    }
  } catch {
    // ignore
  }
  return { text: raw, author: fallbackAuthor }
}

const scrollToBottom = () => {
  nextTick(() => {
    const el = feedRef.value
    if (!el) return
    el.scrollTop = el.scrollHeight
  })
}

const isNearBottom = () => {
  const el = feedRef.value
  if (!el) return true
  const threshold = 40
  return el.scrollHeight - el.scrollTop - el.clientHeight < threshold
}

watch(selectedChannelId, async () => {
  if (!authenticated.value) return
  await loadHistory()
  connectStream()
})

onMounted(() => {
  intervalId = setInterval(pruneWindow, 1000)
  checkSession()
})

onBeforeUnmount(() => {
  disconnectStream()
  if (intervalId) {
    clearInterval(intervalId)
  }
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
  }
})
</script>
