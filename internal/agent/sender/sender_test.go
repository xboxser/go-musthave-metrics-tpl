package sender

import "testing"

func TestSendRequest(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}

	tests := []struct {
		name     string
		args     args
		valError bool
		want     int64
	}{
		{name: "valid parameters counter", args: args{metricType: "counter", metricName: "name", metricValue: "123"}, valError: false, want: int64(123)},
		{name: "valid parameters gauge", args: args{metricType: "gauge", metricName: "2", metricValue: "gfd"}, valError: true, want: int64(1)},
		{name: "invalid metricType", args: args{metricType: "sdfsdfsdfdsfsd", metricName: "2", metricValue: "234"}, valError: true, want: int64(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// todo уточнить как отправлять запрос
			//err := SendRequest(tt.args.metricType, tt.args.metricName, tt.args.metricValue)

			// if (err != nil) != tt.valError {
			// 	t.Errorf("SendRequest: %v", err)
			// 	return
			// }
		})
	}
}
