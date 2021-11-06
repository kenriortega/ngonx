import {
    Box,
    chakra,
    Stack,
    Text,
    VisuallyHidden,
    useColorModeValue,
    Center,
} from '@chakra-ui/react';
import { useStylesApp } from '../hooks/useStyleApp';


export const SocialButton = ({
    children,
    label,
    href,
}) => {
    return (
        <chakra.button
            // bg={useColorModeValue('blackAlpha.100', 'whiteAlpha.100')}
            rounded={'full'}
            w={8}
            h={8}
            cursor={'pointer'}
            as={'a'}
            href={href}
            display={'inline-flex'}
            alignItems={'center'}
            justifyContent={'center'}
            transition={'background 0.3s ease'}
            _hover={{
                // bg: useColorModeValue('blackAlpha.200', 'whiteAlpha.200'),
            }}>
            <VisuallyHidden>{label}</VisuallyHidden>
            {children}
        </chakra.button>
    );
};

const ListHeader = ({ children }) => {
    return (
        <Text fontWeight={'500'} fontSize={'lg'} mb={2}>
            {children}
        </Text>
    );
};

export default function LargeWithNewsletter() {
    const { border, colorBase } = useStylesApp()
    return (
        <Box
            pt="20"
            color={border}>
            <Center>
                <Stack spacing={6}>

                    <Text fontSize={'sm'}>
                        © {new Date().getFullYear()} @kenriortega web page. All rights reserved
                    </Text>
                </Stack>
            </Center>
        </Box>
    );
}