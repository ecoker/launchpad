defmodule TestApp.Chat.MessageTest do
  use ExUnit.Case, async: true

  alias TestApp.Chat.Message

  describe "new/2" do
    test "builds a message with trimmed fields and generated id" do
      msg = Message.new("  alice  ", "  hello world  ")

      assert msg.username == "alice"
      assert msg.body == "hello world"
      assert is_binary(msg.id)
      assert %DateTime{} = msg.sent_at
    end
  end

  describe "display/1" do
    test "formats as 'username: body'" do
      msg = Message.new("bob", "hi there")
      assert Message.display(msg) == "bob: hi there"
    end
  end

  describe "valid?/1" do
    test "returns true for non-empty username and body" do
      assert Message.valid?(Message.new("alice", "hello"))
    end

    test "returns false for empty body" do
      refute Message.valid?(Message.new("alice", "   "))
    end

    test "returns false for empty username" do
      refute Message.valid?(Message.new("   ", "hello"))
    end
  end
end
