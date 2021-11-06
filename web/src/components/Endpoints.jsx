import React, { Fragment, useEffect, useReducer, useState } from 'react'
import EndpointCard from './EndpointCard'
import { Box, Container, Flex, Heading, Text, Center } from "@chakra-ui/react"
import { FaIdCard, FaTable } from 'react-icons/fa'
import EndpointTable from './EndpointTable'
import NoDataImage from './NoData'
import { useQuery } from 'react-query'

const Endpoints = () => {

    const [viewElement, setViewElement] = useState(false)

    const handleViewElement = () => {
        setViewElement(!viewElement)
    }
    const handleSelectedRow = (endpoint) => {
        console.log(endpoint)
    }
    const width = "100%";

    const { isLoading, error, data } = useQuery('endpointsData', () =>
        fetch('/api/v1/mngt/')
            .then(response =>
                response.json()
            )
    )

    if (isLoading) return 'Loading...'

    if (error) return 'An error has occurred: ' + error.message


    return (
        <>
            <Container maxW="container.1sm" px={[0, 4]}>
                <Heading as="h2" m={8} size="md">
                    Ngonx Proxy!!
                </Heading>

                <Heading as="h1" m={8} size="2xl">
                    Service Discovery
                </Heading>
                {

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
                                            data.length > 0
                                                ? <EndpointTable endpoints={data} onSelectedRow={handleSelectedRow} />
                                                : <NoDataImage
                                                    width="300"
                                                    height="300" />
                                        }

                                    </Box>

                                )
                                : (

                                    <div>
                                        {
                                            data.length > 0
                                                ? (<div className={'grid'}>
                                                    {data.map((endpoint) => (
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

                }
            </Container>
        </>
    )
}
export default Endpoints;