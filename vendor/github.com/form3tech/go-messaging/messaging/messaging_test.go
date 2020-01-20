package messaging

import (
	"testing"
)

func Test_SendAndReceiveMessage(t *testing.T) {

	given, when, then := SendAndReceiveTest(t)
	defer then.Teardown()

	given.
		A_message_with_the_following_payload("This is a message!").And().
		A_listener_for_queue_(QueueName)

	when.
		The_message_is_sent_to_queue_(QueueName)

	then.
		The_message_is_received_by_the_listener()
}

func Test_ListenerWithContext_SendAndReceiveMessage(t *testing.T) {
	given, when, then := SendAndReceiveTest(t)
	defer then.Teardown()

	given.
		A_message_with_the_following_payload("This is a message!").And().
		A_listener_with_context_for_queue_(QueueName2)

	when.
		The_message_is_sent_to_queue_(QueueName2)

	then.
		The_message_is_received_by_the_listener().And().
		The_context_must_contains_correlationId()
}

func Test_Notification(t *testing.T) {
	given, when, then := SendAndReceiveTest(t)
	defer then.Teardown()

	given.
		A_message_with_the_following_payload("This is a message!")

	when.
		A_notification_is_sent()

	then.
		The_notification_is_sent_succesfully()
}

func Test_PurgeQueue(t *testing.T) {

	given, when, then := SendAndReceiveTest(t)
	defer then.Teardown()

	given.
		A_message_with_the_following_payload("This is a message!").And().
		The_message_is_sent_to_queue_(QueueName)

	when.
		The_queue_is_purged(QueueName)

	then.
		The_queue_is_empty(QueueName)
}
