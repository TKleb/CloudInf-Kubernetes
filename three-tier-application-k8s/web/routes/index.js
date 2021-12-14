const express = require('express');
const router = express.Router();
const request = require('request');


const api_url = 'http://' + process.env.API_HOST + ':' + process.env.API_PORT + '/time';

/* GET home page. */
router.get('/', function (req, res, next) {
  request(
    {
      method: 'GET',
      url: api_url,
      json: true
    },
    function (error, response, body) {
      if (error || response.statusCode !== 200) {
        console.log("State of the api service was: " + error)
        return res.status(500).send('error running request to ' + api_url);
      } else {
        res.render('index', {
          title: 'Welcome ' + process.env.GROUP_NAME + ' @ ' + process.env.MODULE_NAME,
          request_uuid: body.request_uuid,
          events: body
        });
      }
    }
  );
});

module.exports = router;
