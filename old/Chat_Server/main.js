const Server = require('./server');

/**
 * The graph that makes up the IRC network
 * 
 * The graph layout is an adjacency list of nodes
 * Every server is considered a node in the graph
 * 
 * Every server has a linked list of servers attached to it
 * OR
 * Every server has an array of servers it can contact which represents the connections it has in the graph (i.e. the edges from itself to other nodes (servers))
 * 
 */
const Network = () => {

    const servers = {}; // all nodes in the graph
    const clients = {}; // all clients? idk if this is a good idea or not.

    

    const FindServer = () => {
        // MST algorithm here?
    }

    const AddServer = (server = null) => {
        
    };

    return { AddServer }
}