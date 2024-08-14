import Head from "next/head";
import MapComponent from "./MapComponent";

const Home = () => {
  return (
    <div>
      <Head>
        <title>US Map with Electoral Districts</title>
        <meta
          name="description"
          content="Interactive map of US states and electoral districts"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <h1>US Map with Electoral Districts</h1>
        <MapComponent />
      </main>

      <footer>
        <p>Powered by Next.js and D3.js</p>
      </footer>
    </div>
  );
};

export default Home;
