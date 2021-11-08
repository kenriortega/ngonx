import { useState, useEffect, useMemo, useCallback } from "react";
import { useToast } from "@chakra-ui/react"
export const useSocket = (url) => {
	const toast = useToast()
	const [data, setData] = useState(null);
	const [error, setError] = useState(null);
	const socket = useMemo(
		() => new WebSocket(`ws://0.0.0.0:10001/${url}`),
		[]
	);

	socket.onopen = useCallback(
		(event) => {
			toast({
				title: "Conected",
				description: `ws://0.0.0.0:10001/${url}`,
				status: "success",
				duration: 5000,
				position: "top-right",
				isClosable: true,
			})
			socket.send("endpoints")
		},
		[socket]
	)

	socket.onmessage = useCallback(
		(event) => {
			if (event.data !== null) {
				return setData(JSON.parse(event.data));
			}
		},
		[socket]
	);

	socket.onerror = useCallback(
		(event) => {
			toast({
				title: "Error",
				description: "error",
				status: "error",
				duration: 5000,
				position: "top-right",
				variant: "subtle",
				isClosable: true,
			})
			setError(event);
		},
		[socket]
	);

	socket.onclose = useCallback(
		(event) => {
			toast({
				title: "Closed",
				description: `Closed Socket`,
				status: "error",
				duration: 5000,
				position: "top-right",
				isClosable: true,
			})
		},
		[socket]
	);

	return { data, error };
};
