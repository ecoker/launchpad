defmodule TestAppWeb.ChatLive do
  @moduledoc """
  Real-time chat room with no persistent memory.

  Messages are broadcast via PubSub and held only in each connected
  LiveView's assigns. When the last user disconnects, messages vanish.
  Presence tracks who's currently in the room.
  """

  use TestAppWeb, :live_view

  alias TestApp.Chat.Message
  alias TestAppWeb.Presence

  @topic "chat:lobby"
  @presence_topic "presence:lobby"

  # ── Lifecycle ──────────────────────────────────────────────────────

  @impl true
  def mount(_params, _session, socket) do
    if connected?(socket) do
      Phoenix.PubSub.subscribe(TestApp.PubSub, @topic)
      Phoenix.PubSub.subscribe(TestApp.PubSub, @presence_topic)
    end

    {:ok,
     socket
     |> assign(:messages, [])
     |> assign(:username, "")
     |> assign(:draft, "")
     |> assign(:joined, false)
     |> assign(:points, 0)
     |> assign(:leaderboard, [])}
  end

  # ── Events ─────────────────────────────────────────────────────────

  @impl true
  def handle_event("join", %{"username" => username}, socket) do
    username = String.trim(username)

    if username == "" do
      {:noreply, put_flash(socket, :error, "Username can't be blank")}
    else
      Presence.track(self(), @presence_topic, username, %{points: 0})

      notification = system_message("#{username} joined the room")
      Phoenix.PubSub.broadcast(TestApp.PubSub, @topic, {:new_message, notification})

      {:noreply,
       socket
       |> assign(username: username, joined: true, points: 0)
       |> assign(:leaderboard, build_leaderboard())}
    end
  end

  def handle_event("send", %{"body" => body}, socket) do
    message = Message.new(socket.assigns.username, body)

    if Message.valid?(message) do
      earned = count_letters(body)
      new_points = socket.assigns.points + earned

      Presence.update(self(), @presence_topic, socket.assigns.username, %{points: new_points})
      Phoenix.PubSub.broadcast(TestApp.PubSub, @topic, {:new_message, message})

      {:noreply, assign(socket, draft: "", points: new_points)}
    else
      {:noreply, socket}
    end
  end

  def handle_event("update_draft", %{"body" => body}, socket) do
    {:noreply, assign(socket, draft: body)}
  end

  # ── PubSub / Presence ──────────────────────────────────────────────

  @impl true
  def handle_info({:new_message, message}, socket) do
    {:noreply, assign(socket, messages: socket.assigns.messages ++ [message])}
  end

  def handle_info(%Phoenix.Socket.Broadcast{event: "presence_diff"}, socket) do
    {:noreply, assign(socket, leaderboard: build_leaderboard())}
  end

  # ── Helpers ─────────────────────────────────────────────────────────

  defp build_leaderboard do
    @presence_topic
    |> Presence.list()
    |> Enum.map(fn {name, %{metas: [meta | _]}} -> {name, Map.get(meta, :points, 0)} end)
    |> Enum.sort_by(fn {_name, points} -> points end, :desc)
  end

  defp count_letters(text) do
    text
    |> String.graphemes()
    |> Enum.count(&String.match?(&1, ~r/\p{L}/u))
  end

  defp system_message(text) do
    %Message{
      id: :crypto.strong_rand_bytes(8) |> Base.url_encode64(padding: false),
      username: "system",
      body: text,
      sent_at: DateTime.utc_now()
    }
  end

  # ── Template ───────────────────────────────────────────────────────

  @impl true
  def render(assigns) do
    ~H"""
    <div class="flex flex-col h-[80vh]">
      <h1 class="text-2xl font-bold mb-6 text-center">test-app chat</h1>

      <%= if not @joined do %>
        <.join_form />
      <% else %>
        <div class="flex gap-4 flex-1 min-h-0">
          <div class="flex-1 flex flex-col min-h-0">
            <.message_list messages={@messages} current_user={@username} />
            <.message_form draft={@draft} />
          </div>
          <.leaderboard entries={@leaderboard} />
        </div>
      <% end %>
    </div>
    """
  end

  # ── Function Components ────────────────────────────────────────────

  defp join_form(assigns) do
    ~H"""
    <form phx-submit="join" class="flex gap-3 justify-center mt-12">
      <input
        type="text"
        name="username"
        placeholder="Pick a username…"
        autofocus
        class="bg-gray-800 border border-gray-700 rounded px-4 py-2 text-gray-100
               focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <button
        type="submit"
        class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded font-medium"
      >
        Join
      </button>
    </form>
    """
  end

  attr :messages, :list, required: true
  attr :current_user, :string, required: true

  defp message_list(assigns) do
    ~H"""
    <div
      id="messages"
      class="flex-1 overflow-y-auto space-y-2 mb-4 p-4 bg-gray-800 rounded"
      phx-hook="ScrollBottom"
    >
      <div :for={msg <- @messages} id={msg.id} class="flex flex-col">
        <%= if msg.username == "system" do %>
          <p class="text-yellow-400 text-xs italic text-center py-1">{msg.body}</p>
        <% else %>
          <span class={[
            "text-xs font-semibold",
            if(msg.username == @current_user, do: "text-blue-400", else: "text-green-400")
          ]}>
            {msg.username}
          </span>
          <p class="text-gray-200">{msg.body}</p>
        <% end %>
      </div>
    </div>
    """
  end

  attr :entries, :list, required: true

  defp leaderboard(assigns) do
    ~H"""
    <aside class="w-48 shrink-0 bg-gray-800 rounded p-3 overflow-y-auto">
      <h2 class="text-xs font-bold uppercase text-gray-500 mb-2 tracking-wide">Talky Talk Leaderboard</h2>
      <ol class="space-y-1">
        <li :for={{name, points} <- @entries} class="text-sm text-gray-300 flex items-center justify-between gap-2 min-w-0">
          <span class="flex items-center gap-2 min-w-0">
            <span class="w-2 h-2 rounded-full bg-green-500 inline-block shrink-0"></span>
            <span class="truncate">{name}</span>
          </span>
          <span class="text-xs font-mono text-yellow-400 shrink-0">{points}</span>
        </li>
      </ol>
    </aside>
    """
  end

  attr :draft, :string, required: true

  defp message_form(assigns) do
    ~H"""
    <form phx-submit="send" class="flex gap-3">
      <input
        type="text"
        name="body"
        value={@draft}
        phx-change="update_draft"
        placeholder="Type a message…"
        autofocus
        autocomplete="off"
        class="flex-1 bg-gray-800 border border-gray-700 rounded px-4 py-2 text-gray-100
               focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <button
        type="submit"
        class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded font-medium"
      >
        Send
      </button>
    </form>
    """
  end
end
