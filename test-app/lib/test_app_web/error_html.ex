defmodule TestAppWeb.ErrorHTML do
  @moduledoc "Renders error pages."

  use TestAppWeb, :html

  def render(template, _assigns) do
    Phoenix.Controller.status_message_from_template(template)
  end
end
