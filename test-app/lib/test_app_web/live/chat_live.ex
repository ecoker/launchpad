defmodule TestAppWeb.ChatLive do
  @moduledoc """
  Real-time chat room with no persistent memory.

  Messages are broadcast via PubSub and held only in each connected
  LiveView's assigns. When the last user disconnects, messages vanish.
  """

  use TestAppWeb, :live_view

  alias TestApp.Chat.Message

  @topic "chat:lobby"

  # ── Lifecycle ──────────────────────────────────────────────────────

  @impl true
  def mount(_params, _session, socket) do
    if connected?(socket), do: Phoenix.PubSub.subscribe(TestApp.PubSub, @topic)

    {:ok,
     socket
     |> assign(:messages, [])
     |> assign(:username, "")
     |> assign(:draft, "")
     |> assign(:joined, false)}
  end

  # ── Events ─────────────────────────────────────────────────────────

  @impl true
  def handle_event("join", %{"username" => username}, socket) do
    username = String.trim(username)

    if username == "" do
      {:noreply, put_flash(socket, :error, "Username can't be blank")}
    else
      {:noreply, assign(socket, username: username, joined: true)}
    end
  end

  def handle_event("send", %{"body" => body}, socket) do
    message = Message.new(socket.assigns.username, body)

    if Message.valid?(message) do
      Phoenix.PubSub.broadcast(TestApp.PubSub, @topic, {:new_message, message})
      {:noreply, assign(socket, draft: "")}
    else
      {:noreply, socket}
    end
  end

  def handle_event("update_draft", %{"body" => body}, socket) do
    {:noreply, assign(socket, draft: body)}
  end

  # ── PubSub ─────────────────────────────────────────────────────────

  @impl true
  def handle_info({:new_message, message}, socket) do
    {:noreply, assign(socket, messages: socket.assigns.messages ++ [message])}
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
        <.message_list messages={@messages} current_user={@username} />
        <.message_form draft={@draft} />
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
        <span class={[
          "text-xs font-semibold",
          if(msg.username == @current_user, do: "text-blue-400", else: "text-green-400")
        ]}>
          {msg.username}
        </span>
        <p class="text-gray-200">{msg.body}</p>
      </div>
    </div>
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
