import React from "react";
import {
    Box,
    Flex,
    Text,
    Stack,
    useColorModeValue,
    useColorMode,
    IconButton,
    Link
} from "@chakra-ui/react";
import Logo from "./Logo";
import { FaSun, FaMoon, FaBars } from 'react-icons/fa'
import { useStylesApp } from "../hooks/useStyleApp";
const NavBar = (props) => {
    const [isOpen, setIsOpen] = React.useState(false);
    const toggle = () => setIsOpen(!isOpen);
    const { border } = useStylesApp()
    const { colorMode, toggleColorMode } = useColorMode()

    return (
        <NavBarContainer {...props}>
            <Logo
                w="50px"
            />

            <MenuToggle toggle={toggle} isOpen={isOpen} colorFill={border} />

            <MenuLinks isOpen={isOpen}
                colorMode={colorMode}
                toggleColorMode={toggleColorMode}
            />
        </NavBarContainer>
    );
};

const CloseIcon = ({ fill }) => (
    <svg width="24" viewBox="0 0 18 18" xmlns="http://www.w3.org/2000/svg">
        <title>Close</title>
        <path
            fill={fill}
            d="M9.00023 7.58599L13.9502 2.63599L15.3642 4.04999L10.4142 8.99999L15.3642 13.95L13.9502 15.364L9.00023 10.414L4.05023 15.364L2.63623 13.95L7.58623 8.99999L2.63623 4.04999L4.05023 2.63599L9.00023 7.58599Z"
        />
    </svg>
);

const MenuIcon = ({ fill }) => (
    <svg
        width="24px"
        viewBox="0 0 20 20"
        xmlns="http://www.w3.org/2000/svg"
        fill={fill}
    >
        <title>Menu</title>
        <path d="M0 3h20v2H0V3zm0 6h20v2H0V9zm0 6h20v2H0v-2z" />
    </svg>
);

const MenuToggle = ({ toggle, isOpen, colorFill }) => {
    return (
        <Box display={{ base: "block", md: "none" }} onClick={toggle}>
            {isOpen ? <CloseIcon fill={colorFill} /> : <MenuIcon fill={colorFill} />}
        </Box>
    );
};

const MenuItem = ({ children, isLast = false, to = "/", ...rest }) => {
    return (
        <Link >
            <Text cursor="pointer" display="block" {...rest}>
                {children}
            </Text>
        </Link>
    );
};

const MenuLinks = ({ isOpen, colorMode, toggleColorMode }) => {

    return (
        <Box
            display={{ base: isOpen ? "block" : "none", md: "block" }}
            flexBasis={{ base: "100%", md: "auto" }}
        >
            <Stack
                spacing={8}
                align="center"
                justify={["center", "space-between", "flex-end", "flex-end"]}
                direction={["column", "row", "row", "row"]}
                pt={[4, 4, 0, 0]}
            >
                <MenuItem to="/"
                    className="link-name"
                    textStyle="linkName"
                > ⚙️ Home
                </MenuItem>
                <IconButton
                    aria-label=""
                    mr="4"
                    className="icon-button-name"
                    size="sm"
                    icon={colorMode === 'light' ? <FaSun /> : <FaMoon />}
                    isRound={true}
                    onClick={toggleColorMode}
                    alignSelf="flex-end"
                />
            </Stack>
        </Box >
    );
};

const NavBarContainer = ({ children, ...props }) => {

    return (
        <Flex
            as="nav"
            align="center"
            justify="space-between"
            wrap="wrap"
            w="100%"
            p={8}
            {...props}
        >
            {children}
        </Flex>
    );
};

export default NavBar;