import { Box, Link } from '@chakra-ui/react'
import React from 'react'
import urlParse from 'https://cdn.skypack.dev/url-parse';
const EndpointCard = ({ id, path_url, status }) => {
    const { pathname, hostname, protocol } = new urlParse(path_url)
    return (
        <>
            <span className={"card"} >
                <Box w="500" color={status === "down" ? "red" : "green"}>

                    <h1>{status === "down" ? `☠️ Down` : `✅ UP`}</h1>
                    <h3>Path: {pathname}</h3>
                    <h4>Host: {hostname}</h4>

                </Box>
            </span>

        </>
    )
}

export default EndpointCard