const messagesDiv = document.getElementById('messages');
const startBtn = document.getElementById('startBtn');
let audioContext;
let oscillator;
let gainNode;
let filterNode;
let socket;

startBtn.addEventListener('click', () => {
    if (!audioContext) {
        audioContext = new (window.AudioContext || window.webkitAudioContext)();
    }

    if (audioContext.state === 'suspended') {
        audioContext.resume();
    }

    //TODO: think of something more clever
    socket = new WebSocket('wss://example.com:443/ws');

    socket.addEventListener('open', () => {
        console.log('WebSocket connected');
    });

    socket.addEventListener('close', () => {
        console.log('WebSocket disconnected');
    });

    socket.addEventListener('error', (error) => {
        console.error('WebSocket error:', error);
    });

    // TODO: untangle the spaghetti
    socket.addEventListener('message', (event) => {
        const reader = new FileReader();

        reader.onload = function(e) {
            const arrayBuffer = e.target.result;
            const dataView = new DataView(arrayBuffer);
            const uint16Value = dataView.getUint16(0, false);
            const normalizedFrequency = 50 + (uint16Value / 65535) * (2000 - 50);

            const gain = 1 / Math.sqrt(normalizedFrequency);

            messagesDiv.textContent = `Port: ${uint16Value} - Freq: ${normalizedFrequency.toFixed(2)} Hz - Gain: ${gain.toFixed(2)}`;

            oscillator = audioContext.createOscillator();
            oscillator.type = 'triangle';
            oscillator.frequency.setValueAtTime(normalizedFrequency, audioContext.currentTime);

            gainNode = audioContext.createGain();
            gainNode.gain.setValueAtTime(gain, audioContext.currentTime);
            gainNode.gain.linearRampToValueAtTime(0, audioContext.currentTime + 5);

            filterNode = audioContext.createBiquadFilter();
            filterNode.type = 'lowpass';
            filterNode.frequency.setValueAtTime(300, audioContext.currentTime); // More aggressive filtering

            oscillator.connect(filterNode);
            filterNode.connect(gainNode);
            gainNode.connect(audioContext.destination);

            oscillator.start();
            oscillator.stop(audioContext.currentTime + 5);
        };

        reader.readAsArrayBuffer(event.data);
    });

    messagesDiv.style.display = 'block';
    startBtn.style.display = 'none';
});

window.addEventListener('beforeunload', () => {
    if (socket) {
        socket.close();
    }
});
