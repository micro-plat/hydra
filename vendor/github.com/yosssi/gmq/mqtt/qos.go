package mqtt

// Values of Qos
const (
	// QoS0 represents "QoS 0: At most once delivery".
	QoS0 byte = iota
	// QoS1 represents "QoS 1: At least once delivery".
	QoS1
	// QoS2 represents "QoS 2: Exactly once delivery".
	QoS2
)

// ValidQoS returns true if the input QoS equals to
// QoS0, QoS1 or QoS2.
func ValidQoS(qos byte) bool {
	return qos == QoS0 || qos == QoS1 || qos == QoS2
}
