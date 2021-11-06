import { useColorModeValue, useMediaQuery } from "@chakra-ui/react"

export const useStylesApp = () => {
    const border = useColorModeValue("#282828", "#fff")
    const [isMobileDevice] = useMediaQuery("(max-width: 600px)")

    return {
        border,
        isMobileDevice
    }
}