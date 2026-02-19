defmodule TestAppWeb.Presence do
  @moduledoc """
  Tracks which users are currently in the chat room.

  Built on Phoenix.Presence — backed by a CRDT, no external store needed.
  """

  use Phoenix.Presence,
    otp_app: :test_app,
    pubsub_server: TestApp.PubSub
end
