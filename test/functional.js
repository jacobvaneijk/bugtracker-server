import { check, group, sleep } from 'k6';
import { Rate } from 'k6/metrics';
import http from 'k6/http';

const failureRate = new Rate('check_failure_rate');
const baseUrl = 'http://188.166.121.197:4321'

export let options = {
    stages: [
        {
            duration: '10s',
            target: 50,
        },
    ],
    thresholds: {
        'http_req_duration': [
            'p(95) < 5000',
        ],
        'check_failure_rate': [
            'rate < 0.01',
            {
                abortOnFail: true,
                threshold: 'rate <= 0.05',
            },
        ],
    },
}

export default function () {
    // Test the 'GET /status' endpoint.
    failureRate.add(!check(http.get(`${baseUrl}/status`), {
        'content is present': (response) => response.body.indexOf('API is up and running!') !== -1,
        'http1.1 is used': (response) => response.proto === 'HTTP/1.1',
        'status is 200': (response) => response.status === 200,
    }))

    // Test the 'GET /bugs' endpoint.
    failureRate.add(!check(http.get(`${baseUrl}/bugs?project=1&path=/contact`), {
        'content body is empty array': (response) => response.json().length === 0,
        'status is 200': (response) => response.status === 200,
    }))

    // Test the 'POST /bugs' endpoint.
    const formData = {
        description: 'Description',
        screenshot: '',
        title: 'Title',
        path: '/',
        project_id: 1,
        selection_width: 50,
        selection_height: 150,
        page_width: 1280,
        page_height: 768,
        dot_x: 60,
        dot_y: 160,
    }

    failureRate.add(!check(http.post(`${baseUrl}/bugs`, formData), {
        'status is 200': (response) => response.status === 200,
    }))

    // Sleep for a random amount of time (between 2 and 5 seconds).
    sleep(Math.random() * 3 + 2)
}
