defmodule TestApp.Chat.Message do
  @moduledoc """
  Pure domain representation of a chat message.

  Messages are ephemeral — they exist only in connected processes' memory.
  No persistence, no side effects.
  """

  @type t :: %__MODULE__{
          id: String.t(),
          username: String.t(),
          body: String.t(),
          sent_at: DateTime.t()
        }

  @enforce_keys [:username, :body]
  defstruct [:id, :username, :body, :sent_at]

  @doc "Build a new message from a username and body text."
  @spec new(String.t(), String.t()) :: t()
  def new(username, body) do
    %__MODULE__{
      id: generate_id(),
      username: String.trim(username),
      body: String.trim(body),
      sent_at: DateTime.utc_now()
    }
  end

  @doc "Format the message for display."
  @spec display(t()) :: String.t()
  def display(%__MODULE__{username: username, body: body}) do
    "#{username}: #{body}"
  end

  @doc "Returns true if the message body is not empty after trimming."
  @spec valid?(t()) :: boolean()
  def valid?(%__MODULE__{username: username, body: body}) do
    String.trim(username) != "" and String.trim(body) != ""
  end

  defp generate_id do
    :crypto.strong_rand_bytes(8) |> Base.url_encode64(padding: false)
  end
end
