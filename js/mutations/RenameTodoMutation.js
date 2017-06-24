import {
  commitMutation,
  graphql,
} from 'react-relay';

const mutation = graphql`
  mutation RenameTodoMutation($input: RenameTodoInput!) {
    renameTodo(input:$input) {
      todo {
        id
        text
      }
    }
  }
`;

function getOptimisticResponse(text, todo) {
  return {
    renameTodo: {
      todo: {
        id: todo.id,
        text: text,
      },
    },
  };
}

let tempID = 0;

function commit(
  environment,
  text,
  todo
) {
  return commitMutation(
    environment,
    {
      mutation,
      variables: {
        input: {
          text, id:
          todo.id,
          clientMutationId: tempID++,
        },
      },
      optimisticResponse: getOptimisticResponse(text, todo),
    }
  );
}

export default {commit};
