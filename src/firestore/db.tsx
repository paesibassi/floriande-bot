import { initializeApp } from "firebase/app";
import { getAuth, signInAnonymously } from "firebase/auth";

const firebaseConfig = {
    apiKey: process.env.REACT_APP_APIKEY,
    authDomain: process.env.REACT_APP_AUTHDOMAIN,
    projectId: process.env.REACT_APP_PROJECTID,
};

const app = initializeApp(firebaseConfig);

const authenticateAnonymously = () => {
    return signInAnonymously(getAuth(app));
};

export const authenticate = async () => {
    await authenticateAnonymously();
};
