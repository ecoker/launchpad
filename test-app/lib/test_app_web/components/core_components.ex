defmodule TestAppWeb.CoreComponents do
  @moduledoc """
  Shared function components used across the application.
  """

  use Phoenix.Component

  @doc "Renders a flash message."
  attr :id, :string, default: "flash"
  attr :flash, :map, default: %{}
  attr :kind, :atom, values: [:info, :error]

  def flash_message(assigns) do
    ~H"""
    <div :if={msg = Phoenix.Flash.get(@flash, @kind)} id={@id} role="alert" class="mb-4 p-4 rounded">
      <p>{msg}</p>
    </div>
    """
  end
end
