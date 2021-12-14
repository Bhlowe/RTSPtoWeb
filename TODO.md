# TODO Items

## Questions
What ports are used by the server (webrtc ports, ice/stun, etc.)
Need to list all ports used so that they can be exposed in Dockerfile.

What is your token system (CheckTokenAccess) Need to understand if it is the same as my requirement.

How to automatically determine which streaming type to use for each browser?


## ClientList
Get List of clients viewing streams. It has been partially implemented in bhlowe fork.
Need to get list of HSL and MSE clients

Add ability to delete (disconnect) a viewer by client ID.
When a stream is "finished" or no longer active, the stream list should delete the client record from the list of active clients.

Might make sense to have a "lastTime" timestamp that keeps track of the last time data was sent to a user.
It might be nice to have a bytes sent and start time. Not required. 



## TokenAccess
Need a method of providing security for a camera stream.
I don't know if this is similar to CheckTokenAccess that you started using oauth.
I am open to using oauth, or a simple method of tokens.

Need to be able to allow one user to access 1 or more streams.

The token can be used as the client ID.

struct Token { string token, string streamId }
/token/add?id=xxxx&stream=yyyy  { return success if stream exists }
/token/remove?id=xxxx  { return success/fail } 
/tokens - { return list or map of token:stream pairs } 

Global variable of whether tokens are enabled. 

Create function to check if stream can be created:
createClientID(streamId, token, (oauth info?))
Return clientID if allowed. If tokens are required, ensure token:stream exists.
if token system disabled, and stream exists.





## Access Check on client stream request
Accessing a stream would via http or web socket would need a new parameter




## Other Suggestions
Make all API calls start with /api/ 
Makes routing via reverse proxy easier

