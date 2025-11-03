import PropTypes from 'prop-types';
import Container from './Container.jsx';

/**
 * Semantic section wrapper providing spacing and optional title.
 */
const Section = ({ id, title, description, children, className = '', as = 'section' }) => {
  const Component = as;
  const labelledBy = title ? `${id}-title` : undefined;
  return (
    <Component id={id} aria-labelledby={labelledBy} className={className}>
      <Container className="py-16 sm:py-20">
        {title ? (
          <div className="mb-10 flex flex-col gap-3 text-center">
            {description && (
              <p className="text-sm uppercase tracking-[0.2em] text-ocean/80 dark:text-white/70">{description}</p>
            )}
            <h2 id={labelledBy} className="text-3xl font-bold text-midnight dark:text-white sm:text-4xl">
              {title}
            </h2>
          </div>
        ) : null}
        {children}
      </Container>
    </Component>
  );
};

Section.propTypes = {
  id: PropTypes.string,
  title: PropTypes.string,
  description: PropTypes.string,
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
  as: PropTypes.oneOf(['section', 'div']),
};

export default Section;
