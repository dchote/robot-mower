package kalman

//
// Adapted from https://github.com/TKJElectronics/KalmanFilter
//

import ()

type KalmanFilter struct {
	Q_angle   float64
	Q_bias    float64
	R_measure float64

	angle float64
	bias  float64
	rate  float64

	P [2][2]float64

	Initialized bool
}

func NewKalmanFilter() *KalmanFilter {
	// We will set the variables like so, these can also be tuned by the user
	filter := &KalmanFilter{
		Q_angle:   0.001,
		Q_bias:    0.003,
		R_measure: 0.03,

		angle: 0.0, // Reset the angle
		bias:  0.0, // Reset bias

		P: [2][2]float64{
			{0, 0},
			{0, 0},
		},
		Initialized: false,
	}

	return filter
}

func (k *KalmanFilter) GetAngle(newAngle float64, newRate float64, dt float64) (angle float64) {
	// KasBot V2  -  Kalman filter module - http://www.x-firm.com/?page_id=145
	// Modified by Kristian Lauszus
	// See my blog post for more information: http://blog.tkjelectronics.dk/2012/09/a-practical-approach-to-kalman-filter-and-how-to-implement-it

	// Discrete Kalman filter time update equations - Time Update ("Predict")
	// Update xhat - Project the state ahead
	/* Step 1 */

	k.rate = newRate - k.bias
	k.angle += dt * k.rate

	// Update estimation error covariance - Project the error covariance ahead
	/* Step 2 */
	k.P[0][0] += dt * (dt*k.P[1][1] - k.P[0][1] - k.P[1][0] + k.Q_angle)
	k.P[0][1] -= dt * k.P[1][1]
	k.P[1][0] -= dt * k.P[1][1]
	k.P[1][1] += k.Q_bias * dt

	// Discrete Kalman filter measurement update equations - Measurement Update ("Correct")
	// Calculate Kalman gain - Compute the Kalman gain

	/* Step 4 */
	S := k.P[0][0] + k.R_measure // Estimate error

	/* Step 5 */
	var K [2]float64 // Kalman gain - This is a 2x1 vector
	K[0] = k.P[0][0] / S
	K[1] = k.P[1][0] / S

	// Calculate angle and bias - Update estimate with measurement zk (newAngle)
	/* Step 3 */
	y := newAngle - k.angle // Angle difference

	/* Step 6 */
	k.angle += K[0] * y
	k.bias += K[1] * y

	// Calculate estimation error covariance - Update the error covariance
	/* Step 7 */
	P00_temp := k.P[0][0]
	P01_temp := k.P[0][1]

	k.P[0][0] -= K[0] * P00_temp
	k.P[0][1] -= K[0] * P01_temp
	k.P[1][0] -= K[1] * P00_temp
	k.P[1][1] -= K[1] * P01_temp

	return k.angle
}

func (k *KalmanFilter) SetAngle(newValue float64) {
	k.angle = newValue
	k.Initialized = true
}

func (k *KalmanFilter) GetRate() (value float64) {
	return k.rate
}

func (k *KalmanFilter) SetQangle(newValue float64) {
	k.Q_angle = newValue
}
func (k *KalmanFilter) SetQbias(newValue float64) {
	k.Q_bias = newValue
}
func (k *KalmanFilter) SetRmeasure(newValue float64) {
	k.R_measure = newValue
}

func (k *KalmanFilter) GetQangle() (value float64) {
	return k.Q_angle
}
func (k *KalmanFilter) GetQbias() (value float64) {
	return k.Q_bias
}
func (k *KalmanFilter) GetRmeasure() (value float64) {
	return k.R_measure
}
