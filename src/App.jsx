import Navbar from './components/Navbar.jsx';
import Hero from './components/Hero.jsx';
import Services from './components/Services.jsx';
import Expertises from './components/Expertises.jsx';
import About from './components/About.jsx';
import Contact from './components/Contact.jsx';
import Footer from './components/Footer.jsx';

const App = () => (
  <div className="flex min-h-screen flex-col bg-white text-brand dark:bg-night dark:text-white">
    <Navbar />
    <main>
      <Hero />
      <Services />
      <Expertises />
      <About />
      <Contact />
    </main>
    <Footer />
  </div>
);

export default App;
