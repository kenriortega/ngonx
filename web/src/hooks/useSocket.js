import { useState, useEffect, useMemo, useCallback } from "react";

export const useSocket = (url) => {
	const [data, setData] = useState(null);
	const [error, setError] = useState(null);
	const socket = useMemo(
		() => new WebSocket(`ws://localhost:10001/${url}`),
		[]
	);

	socket.onmessage = useCallback(
		(event) => {
			console.log("Socket message: ", data);
			return setData(event.data);
		},
		[socket]
	);

	socket.onerror = useCallback(
		(error) => {
			console.log("Socket error: ", error);
			setError(error);
		},
		[socket]
	);

	socket.onclose = useCallback(
		(event) => {
			if (event.wasClean) {
				setError(
					`Connection closed cleanly, code=${event.code} reason=${event.reason}`
				);
			} else {
				//  server process killed or network down
				setError("Connection died");
			}
		},
		[socket]
	);

	const closeConnection = useCallback(
		(code, reason) => {
			socket.close(code, reason);
		},
		[socket]
	);

	useEffect(() => {
		return closeConnection(1000, "Clean up");
	}, []);

	return { data, error, close: closeConnection };
};
