const BanList = (nickmask = "") => {

    // TODO
    // maybe pass in name and nickname then create mask here? easier for logging.

    // TODO
    // clean user input nickmask
    
    let banList = {};

    /**
     * Add user to banList
     */
    const Add = () => {
        if (banList.hasOwnProperty(nickmask)) {
            console.log(`${nickmask} is already on the ban list.`);
        } else {
            console.log(`Adding ${nickmask} to the ban list.`);
            blacklist[nickmask] = true;
        }
    }

    /**
     * Remove user from banList
     */
    const Remove = () => {
        if (banList.hasOwnProperty(nickmask)) {
            console.log(`Removing ${nickmask} from the ban list.`);
            delete blacklist[nickmask];
        } else {
            console.log(`${nickmask} is not on the ban list.`);
        }
    }

    return { Add, Remove }
}