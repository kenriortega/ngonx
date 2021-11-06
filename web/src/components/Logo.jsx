import React from "react";
import { Box, Text } from "@chakra-ui/react";

export default function Logo(props) {

    return (
        <Box {...props}>
            <Text
                as="span"
                fontSize="sm"
                fontWeight="bold"
            >ngonx</Text>
        </Box>
    );
}