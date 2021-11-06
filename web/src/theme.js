import { extendTheme } from '@chakra-ui/react'


const theme = extendTheme({
    styles: {
        global: (props) => ({
            "html, body": {
                fontSize: "sm",
                color: props.colorMode === "dark" ? "white" : "spotify.700",
                lineHeight: "tall",
                bg: props.colorMode === "dark" ? "spotify.700" : "white",
            },

        }),
    },
    config: {
        initialColorMode: "system",
        useSystemColorMode: true
    },
    colors: {
        spotify: {
            700: "#282828"
        },

    },
    sizes: {
        container: {
            sectionLanding: "1250px",
            desktop: "1250px",
            desktopInput: "450px",
            "1sm": "700px",
            "2sm": "850px",
            "1xl": "1250px",
            "2xl": "1440px",
        }
    },
    textStyles: {
        filterListItem: {
            fontSize: ['18px', '24px'],
            fontWeight: 'semibold',
            lineHeight: '110%',
            letterSpacing: '-1%',
            paddingLeft: '23px'
        },
        brandName: {
            fontSize: ['18px', '24px'],
            fontWeight: 'semibold',
            lineHeight: '110%',
            letterSpacing: '-1%'
        },
        linkName: {
            fontSize: ['12px', '14px'],
            fontWeight: 'semibold',
            lineHeight: '110%',
            letterSpacing: '-1%',
            justifyContent: "center",
            padding: "4px 16px"
        },
        cardText: {
            fontSize: ['14px', '16px', '18px', '24px'],
            fontWeight: 'semibold',
            lineHeight: '110%',
            letterSpacing: '-1%'
        },
        cardTextSmall: {
            fontSize: ['10px', '12px', '12px', '14px'],
            fontWeight: '500',
            lineHeight: '110%',
            letterSpacing: '-1%'
        },
        cardTextVerySmall: {
            fontSize: ['10px', '10px'],
            fontWeight: '500',
            letterSpacing: '-1%'
        }
    },

})

export default theme