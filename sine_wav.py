import numpy as np
from scipy.io import wavfile

output_name = 'data/sine.wav'
sample_rate = 16000
sin_frequencies = [7999, 4000, 1000, 1]
duration = 1

samples = np.linspace(0, duration, int(sample_rate * duration), endpoint=False)
signal = 0
for freq in sin_frequencies:
    signal += np.sin(2 * np.pi * freq * samples)
signal /= len(sin_frequencies)

wavfile.write(output_name, sample_rate, np.int16(signal * 32767))
