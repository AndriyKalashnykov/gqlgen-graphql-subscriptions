import { gql } from '@apollo/client';

export default gql`
  subscription {
    messageCreated {
      id
      message
    }
  }
`

