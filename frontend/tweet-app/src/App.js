import React, { useState } from 'react';
import { useMutation, useQuery, gql } from '@apollo/client';
import './App.css';

const GET_TWEETS = gql`
  query GetTweets {
    tweets {
      id
      content
    }
  }
`;

const CREATE_TWEET = gql`
  mutation CreateTweet($content: String!) {
    createTweet(content: $content) {
      id
      content
    }
  }
`;

function App() {
  const [tweetContent, setTweetContent] = useState('');
  const [message, setMessage] = useState('');

  const { loading, error, data } = useQuery(GET_TWEETS);

  const [createTweet] = useMutation(CREATE_TWEET, {
    update: (cache, { data: { createTweet } }) => {
      const { tweets } = cache.readQuery({ query: GET_TWEETS });
      cache.writeQuery({
        query: GET_TWEETS,
        data: { tweets: [...tweets, createTweet] },
      });
    },
    onError: (error) => {
      setMessage('An error occurred while creating the tweet.');
    },
  });

  const handleSubmit = async (event) => {
    event.preventDefault();

    try {
      const { data: { createTweet: newTweet } } = await createTweet({
        variables: { content: tweetContent },
      });
      setMessage(`Tweet created: ${newTweet.content}`);
      setTweetContent('');
    } catch (error) {
      console.error('Error:', error);
      setMessage(`An error occurred while creating the tweet: ${error.message}`);
    }
  };

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div className="App">
      <h1>Create a Tweet</h1>
      <form onSubmit={handleSubmit}>
        <label htmlFor="tweetContent">Tweet Content:</label>
        <br />
        <textarea
          id="tweetContent"
          name="tweetContent"
          rows="4"
          cols="50"
          required
          value={tweetContent}
          onChange={(e) => setTweetContent(e.target.value)}
        ></textarea>
        <br />
        <br />
        <input type="submit" value="Create Tweet" />
      </form>

      {message && (
        <div className={message.includes('error') ? 'error-message' : 'success-message'}>
          {message}
        </div>
      )}

      <h2>Tweets</h2>
      <ul>
        {data.tweets.map((tweet) => (
          <li key={tweet.id}>{tweet.content}</li>
        ))}
      </ul>
    </div>
  );
}

export { App };