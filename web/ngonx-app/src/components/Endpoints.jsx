import React, { Fragment, useEffect, useReducer, useState } from 'react'
import EndpointCard from './EndpointCard'
import { Box, Container, Flex, Heading, Text } from "@chakra-ui/react"
import { getAllEndpoints, GetEndpointsReducer, initialStateGetAllEndpoints } from '../lib/getEndpointsReducer'
import { FaIdCard, FaTable } from 'react-icons/fa'
import EndpointTable from './EndpointTable'
import NoDataImage from './NoData'
const Endpoints = () => {
    // const [{ endpoints, loading, errorMessage }, dispatch] =
    //     useReducer(GetEndpointsReducer, initialStateGetAllEndpoints)

    // useEffect(() => {
    //     async function fetchEndpoints() {
    //         try {
    //             let response = await getAllEndpoints(dispatch)
    //             if (!response) return;
    //         } catch (error) {
    //             console.log(error);
    //         }
    //     }

    //     fetchEndpoints()
    // }, [])

    const endpoints = [
        {
            id: "6c4fa765-f6bf-4f00-950f-3315f504cc50",
            path_url: "http://localhost:3000/api/v1/health/",
            status: "down"
        },
        {
            id: "4abbf98f-8a8c-46ca-a2c5-ae8eaefbef60",
            path_url: "http://localhost:3000/api/v1/version/",
            status: "down"
        }
    ]

    const loading = false
    const [viewElement, setViewElement] = useState(false)

    const handleViewElement = () => {
        setViewElement(!viewElement)
    }
    const handleSelectedRow = (poll) => {
        console.log(poll)
    }
    const width = "100%";
    return (
        <>
            <Container maxW="container.1sm" px={[0, 4]}>
                <Heading as="h2" m={8} size="md">
                    Ngonx Proxy!!

                </Heading>
                <Text m={8} size="md">
                    In this board you can find the status check point for all endpoints defined on yml file
                </Text>


                <Heading as="h1" m={8} size="2xl">
                    Service Discovery
                </Heading>
                {
                    loading
                        ?
                        (<Text>Fetching data ...</Text>)
                        :
                        (
                            <>
                                <Box onClick={handleViewElement} cursor="pointer">
                                    {viewElement
                                        ? <FaIdCard size={23} />
                                        : <FaTable size={23} />
                                    }
                                </Box>
                                {
                                    viewElement
                                        ? (
                                            <Box p={23}>

                                                {
                                                    endpoints.length > 0
                                                        ? <EndpointTable endpoints={endpoints} onSelectedRow={handleSelectedRow} />
                                                        : <NoDataImage
                                                            width="300"
                                                            height="300" />
                                                }

                                            </Box>

                                        )
                                        : (

                                            <div>
                                                {
                                                    endpoints.length > 0
                                                        ? (<div className={'grid'}>
                                                            {endpoints.map((endpoint) => (
                                                                <Fragment key={endpoint.id}>

                                                                    <EndpointCard {...endpoint} />
                                                                </Fragment>
                                                            ))}
                                                        </div>)
                                                        : <Center>
                                                            <NoDataImage
                                                                width="300"
                                                                height="300" />
                                                        </Center>
                                                }
                                            </div>
                                        )
                                }
                            </>
                        )

                }
            </Container>
        </>
    )
}
export default Endpoints;