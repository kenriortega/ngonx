import {
	Table,
	Thead,
	Tr,
	Th,
	Tbody,
	Td,
	Tooltip,
	Link,
} from "@chakra-ui/react";
import React from "react";
import urlParse from "https://cdn.skypack.dev/url-parse";

const EndpointTable = ({ endpoints, onSelectedRow }) => {
	return (
		<>
			<Table variant="simple">
				<Thead>
					<Tr>
						{Object.keys(endpoints[0]).map((k) => (
							<Th key={k}>{k}</Th>
						))}
					</Tr>
				</Thead>
				<Tbody>
					{endpoints.map((endpoint) => (
						<Tr onClick={() => onSelectedRow(endpoint)} key={endpoint.id}>
							<Td cursor="pointer" color={endpoint.status === "down" ? "red" : "green"}>
								<Tooltip label={endpoint.path_url} placement="left">
									{endpoint.id.split("-")[0]}
								</Tooltip>
							</Td>

							<Td>{urlParse(endpoint.path_url).pathname}</Td>
							<Td>{endpoint.status === "down" ? "☠️ down" : "✅ up"}</Td>
						</Tr>
					))}
				</Tbody>
			</Table>
		</>
	);
};

export default EndpointTable;
