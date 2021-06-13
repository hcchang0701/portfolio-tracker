import type {HttpFunction} from '@google-cloud/functions-framework/build/src/functions';
import axios from 'axios';
import {createHmac} from 'crypto';

export const listBalance: HttpFunction = async (req, res) => {
  const binance_api_host = 'https://api.binance.com';
  const path = '/api/v3/account';

  const url = binance_api_host + path;
  const binance_api_secret =
    'm4h1k5jNDregkjqa0lW5KmoOKFhK0Hklm6cVJeLNO7B7OcmSSQVnhfgBQmj59O4R';
  const hash = createHmac('sha256', binance_api_secret);

  const headers = {
    'X-MBX-APIKEY':
      'QsoqKGt4xku8xoUCqSZG7YYWeOqjjyUSIdbbLQlGwpZQSvJK6m9bqNkSnOlgbHvj',
    'Content-Type': 'application/json',
  };

  const params: Record<string, number | string> = {
    timestamp: Date.now(),
  };
  const qs = Object.keys(params)
    .map(key => `${key}=${params[key]}`)
    .join('&');
  const signature = hash.update(qs).digest('hex');
  params.signature = signature;

  let result;
  try {
    result = await axios.get(url, {headers, params});
  } catch(error: any) {
    delete error.config;
    result = error;
  }
  
  res.send(result.data);
};
